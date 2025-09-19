package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	SectionServiceItemListSeparator = "<=>"
)

// SectionFilterParams defines parameters for filtering sections
type SectionFilterParams struct {
	// Basic filters
	IDs        []string   `json:"ids"`
	Search     string     `json:"search"`
	SearchMode SearchMode `json:"search_mode"`

	// Content filters
	HasImage   int `json:"has_image"`    // -1 = no image, 0 = all, 1 = has image
	HasBGImage int `json:"has_bg_image"` // -1 = no bg image, 0 = all, 1 = has bg image

	// Count filters
	ParagraphCount int `json:"paragraph_count"` // Minimum number of paragraphs
	ServiceCount   int `json:"service_count"`   // Minimum number of services
	ItemCount      int `json:"item_count"`      // Minimum number of service items

	// Price filters (services)
	MinPrice int `json:"min_price"` // Minimum service price (in cents)
	MaxPrice int `json:"max_price"` // Maximum service price (in cents)

	// Date filters
	CreatedAfter  string `json:"created_after"`  // ISO date string
	CreatedBefore string `json:"created_before"` // ISO date string
	UpdatedAfter  string `json:"updated_after"`  // ISO date string
	UpdatedBefore string `json:"updated_before"` // ISO date string

	// Pagination & Sorting
	Sort  string `json:"sort"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

func NewSectionFilterParamsFromRequest(r *http.Request) SectionFilterParams {
	params := SectionFilterParams{}

	if r.URL.Query().Get("search") != "" {
		params.Search = r.URL.Query().Get("search")
		params.SearchMode = SearchModeFullText
	}

	if r.URL.Query().Get("has_image") != "" {
		params.HasImage, _ = strconv.Atoi(r.URL.Query().Get("has_image"))
	}
	if r.URL.Query().Get("has_bg_image") != "" {
		params.HasBGImage, _ = strconv.Atoi(r.URL.Query().Get("has_bg_image"))
	}

	if r.URL.Query().Get("paragraph_count") != "" {
		params.ParagraphCount, _ = strconv.Atoi(r.URL.Query().Get("paragraph_count"))
	}
	if r.URL.Query().Get("service_count") != "" {
		params.ServiceCount, _ = strconv.Atoi(r.URL.Query().Get("service_count"))
	}
	if r.URL.Query().Get("item_count") != "" {
		params.ItemCount, _ = strconv.Atoi(r.URL.Query().Get("item_count"))
	}

	if r.URL.Query().Get("min_price") != "" {
		params.MinPrice, _ = strconv.Atoi(r.URL.Query().Get("min_price"))
	}
	if r.URL.Query().Get("max_price") != "" {
		params.MaxPrice, _ = strconv.Atoi(r.URL.Query().Get("max_price"))
	}

	if r.URL.Query().Get("created_after") != "" {
		params.CreatedAfter = r.URL.Query().Get("created_after")
	}
	if r.URL.Query().Get("created_before") != "" {
		params.CreatedBefore = r.URL.Query().Get("created_before")
	}
	if r.URL.Query().Get("updated_after") != "" {
		params.UpdatedAfter = r.URL.Query().Get("updated_after")
	}
	if r.URL.Query().Get("updated_before") != "" {
		params.UpdatedBefore = r.URL.Query().Get("updated_before")
	}

	if r.URL.Query().Get("sort") != "" {
		params.Sort = r.URL.Query().Get("sort")
	}
	if r.URL.Query().Get("page") != "" {
		params.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	}
	if r.URL.Query().Get("limit") != "" {
		params.Limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	}

	return params
}

// SectionFilterResult contains filtered sections with pagination info
type SectionFilterResult struct {
	Sections    []*Section `json:"sections"`
	Total       int        `json:"total"`
	Page        int        `json:"page"`
	Limit       int        `json:"limit"`
	TotalPages  int        `json:"total_pages"`
	HasNext     bool       `json:"has_next"`
	HasPrevious bool       `json:"has_previous"`
	HasError    bool       `json:"has_error"`
	Error       string     `json:"error"`
}

var (
	ErrSectionItemIsNotList     = errors.New("content is not a list")
	ErrSectionInsert            = errors.New("failed to insert section")
	ErrSectionParagraphInsert   = errors.New("failed to insert section paragraph")
	ErrSectionServiceInsert     = errors.New("failed to insert section service")
	ErrSectionServiceItemInsert = errors.New("failed to insert section service item")
	ErrSectionNotFound          = errors.New("section not found")
	ErrSectionUpdate            = errors.New("failed to update section")
	ErrSectionParagraphUpdate   = errors.New("failed to update section paragraph")
	ErrSectionServiceUpdate     = errors.New("failed to update section service")
	ErrSectionServiceItemUpdate = errors.New("failed to update section service item")
	ErrSectionDelete            = errors.New("failed to delete section")
)

type Section struct {
	ID string `db:"id" json:"id"`
	// Name is used as a human-readable identifier
	Name string `db:"name" json:"name"`
	// Title is the text shown on the section header
	Title string `db:"title" json:"title"`
	// Image is the URL of the section's image, if any
	// This image is shown beside the content of the section
	Image string `db:"image" json:"image"`
	// BGImage is the URL of the section's background image, if any
	BGImage string `db:"bg_image" json:"bg_image"`

	Paragraphs []SectionParagraph `json:"paragraphs"`
	Services   []SectionService   `json:"services"`

	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type SectionParagraph struct {
	ID        string `db:"id" json:"id"`
	SectionID string `db:"section_id" json:"section_id"`
	// Order is used to ensure paragraphs are shown in the correct order
	//
	// Since order is reserved word in SQL, we use order_idx for the db field name
	Order   int    `db:"order_idx" json:"order"`
	Content string `db:"content" json:"content"`

	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type SectionService struct {
	ID        string `db:"id" json:"id"`
	SectionID string `db:"section_id" json:"section_id"`
	Title     string `db:"title" json:"title"`
	// Price is multiplied by 100 (e.g. 1000 = 10$) and stored as an integer
	//
	// Some services may not have a general price, but their items may have
	// specific prices (e.g. a service with a list of items, each with its own price)
	Price int `db:"price" json:"price"`
	// Description may be empty and a list of items will be shown instead
	// (e.g. a list of included services/products)
	//
	// If description is not empty, it will be shown as a paragraph, and the list of items will be shown below
	Description string `db:"description" json:"description"`

	Items []SectionServiceItem `json:"items"`

	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type SectionServiceItem struct {
	ID        string `db:"id" json:"id"`
	ServiceID string `db:"service_id" json:"service_id"`
	Order     int    `db:"order" json:"order"`
	// Price is multiplied by 100 (e.g. 1000 = 10$) and stored as an integer
	Price   int    `db:"price" json:"price"`
	Content string `db:"content" json:"content"`
	// ContentAsList is flag to indicate if the content should be shown as a list.
	//
	// If true, the content will be split using [SectionServiceItemListSeparator]
	// and populating the ContentList field.
	//
	// If false, the content will be shown as a paragraph.
	ContentAsList bool `db:"content_as_list" json:"content_as_list"`
	// ContentList is a list of items, obtained by separating Content using [SectionServiceItemListSeparator]
	ContentList []string `db:"content_list" json:"content_list"`
	CreatedAt   string   `db:"created_at" json:"created_at"`
	UpdatedAt   string   `db:"updated_at" json:"updated_at"`
}

func (ssi *SectionServiceItem) ParseContentList() error {
	if !ssi.ContentAsList {
		return ErrSectionItemIsNotList
	}

	ssi.ContentList = strings.Split(ssi.Content, SectionServiceItemListSeparator)
	if len(ssi.ContentList) == 1 && ssi.ContentList[0] == "" {
		return ErrSectionItemIsNotList
	}

	return nil
}

func CreateSection(ctx context.Context, section *Section) error {
	conn, err := GetConnWithContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	section.ID = uuid.Must(uuid.NewV7()).String()
	section.CreatedAt = time.Now().Format("2006-01-02T15:04:05-07:00")

	args := pgx.NamedArgs{
		"id":         section.ID,
		"name":       section.Name,
		"title":      section.Title,
		"image":      section.Image,
		"bg_image":   section.BGImage,
		"created_at": section.CreatedAt,
	}
	_, err = tx.Exec(
		ctx,
		`INSERT INTO sections (id, name, title, image, bg_image, created_at) VALUES (@id, @name, @title, @image, @bg_image, @created_at)`,
		args,
	)
	if err != nil {
		return errors.Join(ErrSectionInsert, err)
	}

	for _, paragraph := range section.Paragraphs {
		id := uuid.Must(uuid.NewV7()).String()
		paraArgs := pgx.NamedArgs{
			"id":         id,
			"section_id": section.ID,
			"order":      paragraph.Order,
			"content":    paragraph.Content,
		}
		_, err = tx.Exec(
			ctx,
			`INSERT INTO section_paragraphs (id, section_id, order_idx, content) VALUES (@id, @section_id, @order, @content)`,
			paraArgs,
		)
		if err != nil {
			return errors.Join(ErrSectionParagraphInsert, err)
		}
	}

	for _, service := range section.Services {
		id := uuid.Must(uuid.NewV7()).String()
		serviceArgs := pgx.NamedArgs{
			"id":          id,
			"section_id":  section.ID,
			"title":       service.Title,
			"price":       service.Price,
			"description": service.Description,
		}
		_, err = tx.Exec(
			ctx,
			`INSERT INTO section_service (id, section_id, title, price, description) VALUES (@id, @section_id, @title, @price, @description)`,
			serviceArgs,
		)
		if err != nil {
			return errors.Join(ErrSectionServiceInsert, err)
		}

		for _, item := range service.Items {
			id := uuid.Must(uuid.NewV7()).String()
			itemArgs := pgx.NamedArgs{
				"id":              id,
				"service_id":      service.ID,
				"order":           item.Order,
				"price":           item.Price,
				"content":         item.Content,
				"content_as_list": item.ContentAsList,
			}
			_, err = tx.Exec(
				ctx,
				`INSERT INTO section_service_items (id, service_id, order, price, content, content_as_list) VALUES (@id, @service_id, @order, @price, @content, @content_as_list)`,
				itemArgs,
			)
			if err != nil {
				return errors.Join(ErrSectionServiceItemInsert, err)
			}
		}
	}

	return tx.Commit(ctx)
}

func FindSectionByID(ctx context.Context, id string) (*Section, error) {
	conn, err := GetConnWithContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var section Section
	var paragraphsJSON, servicesJSON []byte
	var (
		sectionName    sql.NullString
		sectionTitle   sql.NullString
		sectionImage   sql.NullString
		sectionBGImage sql.NullString
		sectionCreated sql.NullString
		sectionUpdated sql.NullString
	)

	err = conn.QueryRow(
		ctx,
		`SELECT id, name, title, image, bg_image, created_at, updated_at, paragraphs, services
		 FROM detailed_sections WHERE id = $1`,
		id,
	).Scan(
		&section.ID,
		&sectionName,
		&sectionTitle,
		&sectionImage,
		&sectionBGImage,
		&sectionCreated,
		&sectionUpdated,
		&paragraphsJSON,
		&servicesJSON,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSectionNotFound
		}
		return nil, err
	}

	// Handle nullable string fields
	if sectionName.Valid {
		section.Name = sectionName.String
	}
	if sectionTitle.Valid {
		section.Title = sectionTitle.String
	}
	if sectionImage.Valid {
		section.Image = sectionImage.String
	}
	if sectionBGImage.Valid {
		section.BGImage = sectionBGImage.String
	}
	if sectionCreated.Valid {
		section.CreatedAt = sectionCreated.String
	}
	if sectionUpdated.Valid {
		section.UpdatedAt = sectionUpdated.String
	}

	// Parse paragraphs JSON
	if len(paragraphsJSON) > 0 {
		err = json.Unmarshal(paragraphsJSON, &section.Paragraphs)
		if err != nil {
			return nil, errors.Join(errors.New("failed to unmarshal paragraphs"), err)
		}
	}

	// Parse services JSON
	if len(servicesJSON) > 0 {
		err = json.Unmarshal(servicesJSON, &section.Services)
		if err != nil {
			return nil, errors.Join(errors.New("failed to unmarshal services"), err)
		}
	}

	return &section, nil
}

// UpdateSection updates only the content/display fields of an existing section and its related records.
// Only records with existing IDs will be updated - new records should be created separately.
//
// This function performs change detection to minimize database operations.
func UpdateSection(ctx context.Context, section *Section) error {
	// First, fetch the current section to compare changes
	currentSection, err := FindSectionByID(ctx, section.ID)
	if err != nil {
		return err
	}

	conn, err := GetConnWithContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update main section only if content fields changed
	if sectionChanged(currentSection, section) {
		sectionArgs := pgx.NamedArgs{
			"id":       section.ID,
			"name":     section.Name,
			"title":    section.Title,
			"image":    section.Image,
			"bg_image": section.BGImage,
		}

		_, err = tx.Exec(
			ctx,
			`UPDATE sections SET
				name = @name, title = @title, image = @image, bg_image = @bg_image
			WHERE id = @id`,
			sectionArgs,
		)
		if err != nil {
			return errors.Join(ErrSectionUpdate, err)
		}
	}

	// Update paragraphs that changed
	for _, newPara := range section.Paragraphs {
		currentPara := findParagraphByID(currentSection.Paragraphs, newPara.ID)
		if currentPara == nil {
			// Skip paragraphs that don't exist - updates should only modify existing records
			continue
		}
		if paragraphChanged(currentPara, &newPara) {
			paraArgs := pgx.NamedArgs{
				"id":      newPara.ID,
				"order":   newPara.Order,
				"content": newPara.Content,
			}
			_, err = tx.Exec(
				ctx,
				`UPDATE section_paragraphs SET
					order_idx = @order, content = @content
				WHERE id = @id`,
				paraArgs,
			)
			if err != nil {
				return errors.Join(ErrSectionParagraphUpdate, err)
			}
		}
	}

	// Update services that changed
	for _, newService := range section.Services {
		currentService := findServiceByID(currentSection.Services, newService.ID)
		if currentService == nil {
			// Skip services that don't exist - updates should only modify existing records
			continue
		}
		if serviceChanged(currentService, &newService) {
			serviceArgs := pgx.NamedArgs{
				"id":          newService.ID,
				"title":       newService.Title,
				"price":       newService.Price,
				"description": newService.Description,
			}
			_, err = tx.Exec(
				ctx,
				`UPDATE section_service SET
					title = @title, price = @price, description = @description
				WHERE id = @id`,
				serviceArgs,
			)
			if err != nil {
				return errors.Join(ErrSectionServiceUpdate, err)
			}
		}

		// Update service items that changed
		for _, newItem := range newService.Items {
			currentItem := findServiceItemByID(currentService.Items, newItem.ID)
			if currentItem == nil {
				// Skip items that don't exist - updates should only modify existing records
				continue
			}
			if serviceItemChanged(currentItem, &newItem) {
				itemArgs := pgx.NamedArgs{
					"id":              newItem.ID,
					"order":           newItem.Order,
					"price":           newItem.Price,
					"content":         newItem.Content,
					"content_as_list": newItem.ContentAsList,
				}
				_, err = tx.Exec(
					ctx,
					`UPDATE section_service_items SET
						order_idx = @order, price = @price, content = @content, content_as_list = @content_as_list
					WHERE id = @id`,
					itemArgs,
				)
				if err != nil {
					return errors.Join(ErrSectionServiceItemUpdate, err)
				}
			}
		}
	}

	return tx.Commit(ctx)
}

// Helper functions for change detection
func sectionChanged(current, updated *Section) bool {
	return current.Name != updated.Name ||
		current.Title != updated.Title ||
		current.Image != updated.Image ||
		current.BGImage != updated.BGImage
}

func paragraphChanged(current, updated *SectionParagraph) bool {
	return current.Order != updated.Order ||
		current.Content != updated.Content
}

func serviceChanged(current, updated *SectionService) bool {
	return current.Title != updated.Title ||
		current.Price != updated.Price ||
		current.Description != updated.Description
}

func serviceItemChanged(current, updated *SectionServiceItem) bool {
	return current.Order != updated.Order ||
		current.Price != updated.Price ||
		current.Content != updated.Content ||
		current.ContentAsList != updated.ContentAsList
}

// Helper functions to find items by ID
func findParagraphByID(paragraphs []SectionParagraph, id string) *SectionParagraph {
	for _, p := range paragraphs {
		if p.ID == id {
			return &p
		}
	}
	return nil
}

func findServiceByID(services []SectionService, id string) *SectionService {
	for _, s := range services {
		if s.ID == id {
			return &s
		}
	}
	return nil
}

func findServiceItemByID(items []SectionServiceItem, id string) *SectionServiceItem {
	for _, i := range items {
		if i.ID == id {
			return &i
		}
	}
	return nil
}

// DeleteSection removes a section and all its related records (paragraphs, services, service items)
// from the database. The foreign key constraints with ON DELETE CASCADE ensure all related
// records are automatically removed when the main section is deleted.
func DeleteSection(ctx context.Context, id string) error {
	conn, err := GetConnWithContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Verify the section exists before attempting deletion
	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sections WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return errors.Join(ErrSectionDelete, err)
	}
	if !exists {
		return ErrSectionNotFound
	}

	// Delete the main section record
	// Due to ON DELETE CASCADE constraints, this will automatically delete:
	// - All section_paragraphs records
	// - All section_service records
	// - All section_service_items records (via the service deletion cascade)
	_, err = tx.Exec(ctx, "DELETE FROM sections WHERE id = $1", id)
	if err != nil {
		return errors.Join(ErrSectionDelete, err)
	}

	return tx.Commit(ctx)
}

// FindAllSections retrieves all sections with their complete details (paragraphs, services, service items)
// from the database. This function uses the detailed_sections view for optimal performance.
func FindAllSections(ctx context.Context) ([]*Section, error) {
	conn, err := GetConnWithContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`SELECT id, name, title, image, bg_image, created_at, updated_at, paragraphs, services
		 FROM detailed_sections
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []*Section

	for rows.Next() {
		var section Section
		var paragraphsJSON, servicesJSON []byte
		var (
			sectionName    sql.NullString
			sectionTitle   sql.NullString
			sectionImage   sql.NullString
			sectionBGImage sql.NullString
			sectionCreated sql.NullString
			sectionUpdated sql.NullString
		)

		err = rows.Scan(
			&section.ID,
			&sectionName,
			&sectionTitle,
			&sectionImage,
			&sectionBGImage,
			&sectionCreated,
			&sectionUpdated,
			&paragraphsJSON,
			&servicesJSON,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable string fields
		if sectionName.Valid {
			section.Name = sectionName.String
		}
		if sectionTitle.Valid {
			section.Title = sectionTitle.String
		}
		if sectionImage.Valid {
			section.Image = sectionImage.String
		}
		if sectionBGImage.Valid {
			section.BGImage = sectionBGImage.String
		}
		if sectionCreated.Valid {
			section.CreatedAt = sectionCreated.String
		}
		if sectionUpdated.Valid {
			section.UpdatedAt = sectionUpdated.String
		}

		// Parse paragraphs JSON
		if len(paragraphsJSON) > 0 {
			err = json.Unmarshal(paragraphsJSON, &section.Paragraphs)
			if err != nil {
				return nil, errors.Join(errors.New("failed to unmarshal paragraphs"), err)
			}
		}

		// Parse services JSON
		if len(servicesJSON) > 0 {
			err = json.Unmarshal(servicesJSON, &section.Services)
			if err != nil {
				return nil, errors.Join(errors.New("failed to unmarshal services"), err)
			}
		}

		sections = append(sections, &section)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sections, nil
}

// FilterSections filters sections based on provided parameters with pagination and sorting
func FilterSections(ctx context.Context, filters SectionFilterParams) (*SectionFilterResult, error) {
	conn, err := GetConnWithContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	// Set defaults
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 || filters.Limit > 100 {
		filters.Limit = 20
	}
	if filters.SearchMode == "" {
		filters.SearchMode = SearchModeFullText
	}

	// Build query conditions and named arguments
	conditions, namedArgs := buildSectionQueryConditions(filters)

	// Base query using detailed_sections view
	baseQuery := "FROM detailed_sections"
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Get total count
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err = conn.QueryRow(ctx, countQuery, namedArgs).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Calculate pagination
	offset := (filters.Page - 1) * filters.Limit
	totalPages := int(math.Ceil(float64(total) / float64(filters.Limit)))

	// Add pagination to named args
	namedArgs["limit"] = filters.Limit
	namedArgs["offset"] = offset

	// Build main query
	selectClause := `SELECT id, name, title, image, bg_image, created_at, updated_at, paragraphs, services, search_vector`
	searchRankClause := buildSectionSearchRankSelect(filters)
	if searchRankClause != "" {
		selectClause += ", " + searchRankClause
	}

	orderByClause := buildSectionOrderByClause(filters)

	mainQuery := fmt.Sprintf("%s %s %s LIMIT @limit OFFSET @offset",
		selectClause, baseQuery, orderByClause)

	// Execute main query
	rows, err := conn.Query(ctx, mainQuery, namedArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to execute main query: %w", err)
	}
	defer rows.Close()

	// Parse results
	var sections []*Section
	for rows.Next() {
		var section Section
		var paragraphsJSON, servicesJSON []byte
		var searchVector []byte // tsvector as byte array
		var (
			sectionName    sql.NullString
			sectionTitle   sql.NullString
			sectionImage   sql.NullString
			sectionBGImage sql.NullString
			sectionCreated sql.NullString
			sectionUpdated sql.NullString
		)

		err = rows.Scan(
			&section.ID,
			&sectionName,
			&sectionTitle,
			&sectionImage,
			&sectionBGImage,
			&sectionCreated,
			&sectionUpdated,
			&paragraphsJSON,
			&servicesJSON,
			&searchVector,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan section: %w", err)
		}

		// Handle nullable string fields
		if sectionName.Valid {
			section.Name = sectionName.String
		}
		if sectionTitle.Valid {
			section.Title = sectionTitle.String
		}
		if sectionImage.Valid {
			section.Image = sectionImage.String
		}
		if sectionBGImage.Valid {
			section.BGImage = sectionBGImage.String
		}
		if sectionCreated.Valid {
			section.CreatedAt = sectionCreated.String
		}
		if sectionUpdated.Valid {
			section.UpdatedAt = sectionUpdated.String
		}

		// Parse paragraphs JSON
		if len(paragraphsJSON) > 0 {
			err = json.Unmarshal(paragraphsJSON, &section.Paragraphs)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal paragraphs: %w", err)
			}
		}

		// Parse services JSON
		if len(servicesJSON) > 0 {
			err = json.Unmarshal(servicesJSON, &section.Services)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal services: %w", err)
			}
		}

		sections = append(sections, &section)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Build result
	result := &SectionFilterResult{
		Sections:    sections,
		Total:       total,
		Page:        filters.Page,
		Limit:       filters.Limit,
		TotalPages:  totalPages,
		HasNext:     filters.Page < totalPages,
		HasPrevious: filters.Page > 1,
		HasError:    false,
		Error:       "",
	}

	return result, nil
}

// buildSectionQueryConditions builds WHERE conditions and named arguments for section filtering
func buildSectionQueryConditions(filters SectionFilterParams) ([]string, pgx.NamedArgs) {
	var conditions []string
	namedArgs := make(pgx.NamedArgs)

	// ID filter
	if len(filters.IDs) > 0 {
		conditions = append(conditions, "id = ANY(@ids)")
		namedArgs["ids"] = filters.IDs
	}

	// Search condition
	if filters.Search != "" {
		switch filters.SearchMode {
		case SearchModeFullText:
			// True full-text search using search vector
			conditions = append(conditions, "search_vector @@ plainto_tsquery('spanish', @search_query)")
			namedArgs["search_query"] = filters.Search

		case SearchModeExact:
			// Exact match search
			conditions = append(conditions, `
				(name ILIKE @exact_search OR
				 title ILIKE @exact_search OR
				 paragraphs::text ILIKE @exact_search OR
				 services::text ILIKE @exact_search)`)
			namedArgs["exact_search"] = filters.Search

		case SearchModeFuzzy:
			// Fuzzy search (LIKE with wildcards)
			conditions = append(conditions, `
				(name ILIKE @fuzzy_search OR
				 title ILIKE @fuzzy_search OR
				 paragraphs::text ILIKE @fuzzy_search OR
				 services::text ILIKE @fuzzy_search)`)
			namedArgs["fuzzy_search"] = "%" + filters.Search + "%"
		}
	}

	// Image filters
	if filters.HasImage > 0 {
		conditions = append(conditions, "image IS NOT NULL AND image != ''")
	} else if filters.HasImage < 0 {
		conditions = append(conditions, "image IS NULL OR image = ''")
	}

	if filters.HasBGImage > 0 {
		conditions = append(conditions, "bg_image IS NOT NULL AND bg_image != ''")
	} else if filters.HasBGImage < 0 {
		conditions = append(conditions, "bg_image IS NULL OR bg_image = ''")
	}

	// Count filters
	if filters.ParagraphCount > 0 {
		conditions = append(conditions, "json_array_length(paragraphs) >= @paragraph_count")
		namedArgs["paragraph_count"] = filters.ParagraphCount
	}

	if filters.ServiceCount > 0 {
		conditions = append(conditions, "json_array_length(services) >= @service_count")
		namedArgs["service_count"] = filters.ServiceCount
	}

	if filters.ItemCount > 0 {
		conditions = append(conditions, `
			(SELECT COALESCE(SUM(json_array_length(svc->'items')), 0)
			 FROM json_array_elements(services) AS svc) >= @item_count`)
		namedArgs["item_count"] = filters.ItemCount
	}

	// Price filters
	if filters.MinPrice > 0 || filters.MaxPrice > 0 {
		if filters.MinPrice > 0 && filters.MaxPrice > 0 {
			conditions = append(conditions, `
				EXISTS (SELECT 1 FROM json_array_elements(services) AS svc
				        WHERE (svc->>'price')::int BETWEEN @min_price AND @max_price)`)
			namedArgs["min_price"] = filters.MinPrice
			namedArgs["max_price"] = filters.MaxPrice
		} else if filters.MinPrice > 0 {
			conditions = append(conditions, `
				EXISTS (SELECT 1 FROM json_array_elements(services) AS svc
				        WHERE (svc->>'price')::int >= @min_price)`)
			namedArgs["min_price"] = filters.MinPrice
		} else if filters.MaxPrice > 0 {
			conditions = append(conditions, `
				EXISTS (SELECT 1 FROM json_array_elements(services) AS svc
				        WHERE (svc->>'price')::int <= @max_price)`)
			namedArgs["max_price"] = filters.MaxPrice
		}
	}

	// Date filters
	if filters.CreatedAfter != "" {
		conditions = append(conditions, "created_at >= @created_after")
		namedArgs["created_after"] = filters.CreatedAfter
	}
	if filters.CreatedBefore != "" {
		conditions = append(conditions, "created_at <= @created_before")
		namedArgs["created_before"] = filters.CreatedBefore
	}
	if filters.UpdatedAfter != "" {
		conditions = append(conditions, "updated_at >= @updated_after")
		namedArgs["updated_after"] = filters.UpdatedAfter
	}
	if filters.UpdatedBefore != "" {
		conditions = append(conditions, "updated_at <= @updated_before")
		namedArgs["updated_before"] = filters.UpdatedBefore
	}

	return conditions, namedArgs
}

// buildSectionSearchRankSelect adds search ranking when using full-text search
func buildSectionSearchRankSelect(filters SectionFilterParams) string {
	if filters.Search != "" && filters.SearchMode == SearchModeFullText {
		return "ts_rank(search_vector, plainto_tsquery('spanish', @search_query)) as search_rank"
	}
	return ""
}

// buildSectionOrderByClause builds the ORDER BY clause for section queries
func buildSectionOrderByClause(filters SectionFilterParams) string {
	// If using full-text search with a query, prioritize search ranking
	if filters.Search != "" && filters.SearchMode == SearchModeFullText {
		switch strings.ToLower(filters.Sort) {
		case "relevance", "":
			return "ORDER BY search_rank DESC, name ASC"
		case "name_asc", "name":
			return "ORDER BY name ASC"
		case "name_desc":
			return "ORDER BY name DESC"
		case "title_asc":
			return "ORDER BY title ASC, search_rank DESC"
		case "title_desc":
			return "ORDER BY title DESC, search_rank DESC"
		case "created_asc":
			return "ORDER BY created_at ASC, search_rank DESC"
		case "created_desc", "newest":
			return "ORDER BY created_at DESC, search_rank DESC"
		case "updated_asc":
			return "ORDER BY updated_at ASC, search_rank DESC"
		case "updated_desc", "recent":
			return "ORDER BY updated_at DESC, search_rank DESC"
		default:
			return "ORDER BY search_rank DESC, name ASC"
		}
	}

	// Regular sorting without search ranking
	switch strings.ToLower(filters.Sort) {
	case "name_asc", "name", "":
		return "ORDER BY name ASC"
	case "name_desc":
		return "ORDER BY name DESC"
	case "title_asc":
		return "ORDER BY title ASC"
	case "title_desc":
		return "ORDER BY title DESC"
	case "created_asc":
		return "ORDER BY created_at ASC"
	case "created_desc", "newest":
		return "ORDER BY created_at DESC"
	case "updated_asc":
		return "ORDER BY updated_at ASC"
	case "updated_desc", "recent":
		return "ORDER BY updated_at DESC"
	case "paragraphs_asc":
		return "ORDER BY json_array_length(paragraphs) ASC, name ASC"
	case "paragraphs_desc":
		return "ORDER BY json_array_length(paragraphs) DESC, name ASC"
	case "services_asc":
		return "ORDER BY json_array_length(services) ASC, name ASC"
	case "services_desc":
		return "ORDER BY json_array_length(services) DESC, name ASC"
	default:
		return "ORDER BY name ASC"
	}
}

// AddServiceToSection adds a new service to an existing section
// This function creates the service and any associated service items in a transaction
func AddServiceToSection(ctx context.Context, sectionID string, service *SectionService) error {
	conn, err := GetConnWithContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Verify the section exists
	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sections WHERE id = $1)", sectionID).Scan(&exists)
	if err != nil {
		return errors.Join(errors.New("failed to verify section existence"), err)
	}
	if !exists {
		return ErrSectionNotFound
	}

	// Generate UUID for the service
	service.ID = uuid.Must(uuid.NewV7()).String()
	service.SectionID = sectionID
	service.CreatedAt = time.Now().Format("2006-01-02T15:04:05-07:00")

	// Insert the service
	serviceArgs := pgx.NamedArgs{
		"id":          service.ID,
		"section_id":  service.SectionID,
		"title":       service.Title,
		"price":       service.Price,
		"description": service.Description,
		"created_at":  service.CreatedAt,
	}
	_, err = tx.Exec(
		ctx,
		`INSERT INTO section_service (id, section_id, title, price, description, created_at)
		 VALUES (@id, @section_id, @title, @price, @description, @created_at)`,
		serviceArgs,
	)
	if err != nil {
		return errors.Join(ErrSectionServiceInsert, err)
	}

	// Insert service items if any
	for _, item := range service.Items {
		item.ID = uuid.Must(uuid.NewV7()).String()
		item.ServiceID = service.ID
		item.CreatedAt = time.Now().Format("2006-01-02T15:04:05-07:00")

		itemArgs := pgx.NamedArgs{
			"id":              item.ID,
			"service_id":      item.ServiceID,
			"order":           item.Order,
			"price":           item.Price,
			"content":         item.Content,
			"content_as_list": item.ContentAsList,
			"created_at":      item.CreatedAt,
		}
		_, err = tx.Exec(
			ctx,
			`INSERT INTO section_service_items (id, service_id, order_idx, price, content, content_as_list, created_at)
			 VALUES (@id, @service_id, @order, @price, @content, @content_as_list, @created_at)`,
			itemArgs,
		)
		if err != nil {
			return errors.Join(ErrSectionServiceItemInsert, err)
		}
	}

	return tx.Commit(ctx)
}

// UpdateSectionWithAdditions updates an existing section and can add new paragraphs, services, and service items
// This function handles both updates to existing records and additions of new records
//
// For records with existing IDs: updates them if they changed
// For records without IDs: creates new records with generated UUIDs
//
// This provides a comprehensive update mechanism that can handle complex section modifications
func UpdateSectionWithAdditions(ctx context.Context, section *Section) error {
	// First, fetch the current section to compare changes
	currentSection, err := FindSectionByID(ctx, section.ID)
	if err != nil {
		return err
	}

	conn, err := GetConnWithContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update main section only if content fields changed
	if sectionChanged(currentSection, section) {
		sectionArgs := pgx.NamedArgs{
			"id":       section.ID,
			"name":     section.Name,
			"title":    section.Title,
			"image":    section.Image,
			"bg_image": section.BGImage,
		}

		_, err = tx.Exec(
			ctx,
			`UPDATE sections SET
				name = @name, title = @title, image = @image, bg_image = @bg_image
			WHERE id = @id`,
			sectionArgs,
		)
		if err != nil {
			return errors.Join(ErrSectionUpdate, err)
		}
	}

	// Handle paragraphs: update existing ones and add new ones
	for _, newPara := range section.Paragraphs {
		if newPara.ID != "" {
			// Existing paragraph - update if changed
			currentPara := findParagraphByID(currentSection.Paragraphs, newPara.ID)
			if currentPara == nil {
				// Skip paragraphs that don't exist in current section
				continue
			}
			if paragraphChanged(currentPara, &newPara) {
				paraArgs := pgx.NamedArgs{
					"id":      newPara.ID,
					"order":   newPara.Order,
					"content": newPara.Content,
				}
				_, err = tx.Exec(
					ctx,
					`UPDATE section_paragraphs SET
						order_idx = @order, content = @content
					WHERE id = @id`,
					paraArgs,
				)
				if err != nil {
					return errors.Join(ErrSectionParagraphUpdate, err)
				}
			}
		} else {
			// New paragraph - insert
			newPara.ID = uuid.Must(uuid.NewV7()).String()
			newPara.SectionID = section.ID
			newPara.CreatedAt = time.Now().Format("2006-01-02T15:04:05-07:00")

			paraArgs := pgx.NamedArgs{
				"id":         newPara.ID,
				"section_id": newPara.SectionID,
				"order":      newPara.Order,
				"content":    newPara.Content,
				"created_at": newPara.CreatedAt,
			}
			_, err = tx.Exec(
				ctx,
				`INSERT INTO section_paragraphs (id, section_id, order_idx, content, created_at)
				 VALUES (@id, @section_id, @order, @content, @created_at)`,
				paraArgs,
			)
			if err != nil {
				return errors.Join(ErrSectionParagraphInsert, err)
			}
		}
	}

	// Handle services: update existing ones and add new ones
	for _, newService := range section.Services {
		if newService.ID != "" {
			// Existing service - update if changed
			currentService := findServiceByID(currentSection.Services, newService.ID)
			if currentService == nil {
				// Skip services that don't exist in current section
				continue
			}
			if serviceChanged(currentService, &newService) {
				serviceArgs := pgx.NamedArgs{
					"id":          newService.ID,
					"title":       newService.Title,
					"price":       newService.Price,
					"description": newService.Description,
				}
				_, err = tx.Exec(
					ctx,
					`UPDATE section_service SET
						title = @title, price = @price, description = @description
					WHERE id = @id`,
					serviceArgs,
				)
				if err != nil {
					return errors.Join(ErrSectionServiceUpdate, err)
				}
			}

			// Handle service items for existing service
			for _, newItem := range newService.Items {
				if newItem.ID != "" {
					// Existing service item - update if changed
					currentItem := findServiceItemByID(currentService.Items, newItem.ID)
					if currentItem == nil {
						// Skip items that don't exist in current service
						continue
					}
					if serviceItemChanged(currentItem, &newItem) {
						itemArgs := pgx.NamedArgs{
							"id":              newItem.ID,
							"order":           newItem.Order,
							"price":           newItem.Price,
							"content":         newItem.Content,
							"content_as_list": newItem.ContentAsList,
						}
						_, err = tx.Exec(
							ctx,
							`UPDATE section_service_items SET
								order_idx = @order, price = @price, content = @content, content_as_list = @content_as_list
							WHERE id = @id`,
							itemArgs,
						)
						if err != nil {
							return errors.Join(ErrSectionServiceItemUpdate, err)
						}
					}
				} else {
					// New service item - insert
					newItem.ID = uuid.Must(uuid.NewV7()).String()
					newItem.ServiceID = newService.ID
					newItem.CreatedAt = time.Now().Format("2006-01-02T15:04:05-07:00")

					itemArgs := pgx.NamedArgs{
						"id":              newItem.ID,
						"service_id":      newItem.ServiceID,
						"order":           newItem.Order,
						"price":           newItem.Price,
						"content":         newItem.Content,
						"content_as_list": newItem.ContentAsList,
						"created_at":      newItem.CreatedAt,
					}
					_, err = tx.Exec(
						ctx,
						`INSERT INTO section_service_items (id, service_id, order_idx, price, content, content_as_list, created_at)
						 VALUES (@id, @service_id, @order, @price, @content, @content_as_list, @created_at)`,
						itemArgs,
					)
					if err != nil {
						return errors.Join(ErrSectionServiceItemInsert, err)
					}
				}
			}
		} else {
			// New service - insert
			newService.ID = uuid.Must(uuid.NewV7()).String()
			newService.SectionID = section.ID
			newService.CreatedAt = time.Now().Format("2006-01-02T15:04:05-07:00")

			serviceArgs := pgx.NamedArgs{
				"id":          newService.ID,
				"section_id":  newService.SectionID,
				"title":       newService.Title,
				"price":       newService.Price,
				"description": newService.Description,
				"created_at":  newService.CreatedAt,
			}
			_, err = tx.Exec(
				ctx,
				`INSERT INTO section_service (id, section_id, title, price, description, created_at)
				 VALUES (@id, @section_id, @title, @price, @description, @created_at)`,
				serviceArgs,
			)
			if err != nil {
				return errors.Join(ErrSectionServiceInsert, err)
			}

			// Handle service items for new service
			for _, newItem := range newService.Items {
				newItem.ID = uuid.Must(uuid.NewV7()).String()
				newItem.ServiceID = newService.ID
				newItem.CreatedAt = time.Now().Format("2006-01-02T15:04:05-07:00")

				itemArgs := pgx.NamedArgs{
					"id":              newItem.ID,
					"service_id":      newItem.ServiceID,
					"order":           newItem.Order,
					"price":           newItem.Price,
					"content":         newItem.Content,
					"content_as_list": newItem.ContentAsList,
					"created_at":      newItem.CreatedAt,
				}
				_, err = tx.Exec(
					ctx,
					`INSERT INTO section_service_items (id, service_id, order_idx, price, content, content_as_list, created_at)
					 VALUES (@id, @service_id, @order, @price, @content, @content_as_list, @created_at)`,
					itemArgs,
				)
				if err != nil {
					return errors.Join(ErrSectionServiceItemInsert, err)
				}
			}
		}
	}

	return tx.Commit(ctx)
}

// AddParagraphToSection adds a new paragraph to an existing section
// This function creates the paragraph with the specified order and content
func AddParagraphToSection(ctx context.Context, sectionID string, paragraph *SectionParagraph) error {
	conn, err := GetConnWithContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Verify the section exists
	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sections WHERE id = $1)", sectionID).Scan(&exists)
	if err != nil {
		return errors.Join(errors.New("failed to verify section existence"), err)
	}
	if !exists {
		return ErrSectionNotFound
	}

	// Generate UUID for the paragraph
	paragraph.ID = uuid.Must(uuid.NewV7()).String()
	paragraph.SectionID = sectionID
	paragraph.CreatedAt = time.Now().Format("2006-01-02T15:04:05-07:00")

	// Insert the paragraph
	paraArgs := pgx.NamedArgs{
		"id":         paragraph.ID,
		"section_id": paragraph.SectionID,
		"order":      paragraph.Order,
		"content":    paragraph.Content,
		"created_at": paragraph.CreatedAt,
	}
	_, err = tx.Exec(
		ctx,
		`INSERT INTO section_paragraphs (id, section_id, order_idx, content, created_at)
		 VALUES (@id, @section_id, @order, @content, @created_at)`,
		paraArgs,
	)
	if err != nil {
		return errors.Join(ErrSectionParagraphInsert, err)
	}

	return tx.Commit(ctx)
}

// AddMultipleParagraphsToSection adds multiple paragraphs to an existing section
// This is a convenience function that calls AddParagraphToSection for each paragraph
func AddMultipleParagraphsToSection(ctx context.Context, sectionID string, paragraphs []*SectionParagraph) error {
	for _, paragraph := range paragraphs {
		if err := AddParagraphToSection(ctx, sectionID, paragraph); err != nil {
			return err
		}
	}
	return nil
}

// AddMultipleServicesToSection adds multiple services to an existing section
// This is a convenience function that calls AddServiceToSection for each service
func AddMultipleServicesToSection(ctx context.Context, sectionID string, services []*SectionService) error {
	for _, service := range services {
		if err := AddServiceToSection(ctx, sectionID, service); err != nil {
			return err
		}
	}
	return nil
}

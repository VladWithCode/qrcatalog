package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vladwithcode/qrcatalog/internal"
)

type Subcategory struct {
	ID              string `db:"id" json:"id"`
	Name            string `db:"name" json:"name"`
	Slug            string `db:"slug" json:"slug"`
	Description     string `db:"description" json:"description"`
	LongDescription string `db:"long_description" json:"longDescription"`
	DisplayImg      string `db:"display_img" json:"displayImg"`
	DisplayImgID    string `db:"display_img_id" json:"displayImgId"`
	CategoryID      string `db:"category_id" json:"categoryId"`
	CategoryName    string `db:"category_name" json:"categoryName"`
	ProductCount    int    `db:"product_count" json:"productCount"`
}

func CreateSubcategory(subcategory *Subcategory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := GetConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	id, err := uuid.NewV7()
	if err != nil {
		return ErrUUIDFail
	}

	displayImg := sql.NullString{
		String: subcategory.DisplayImg,
		Valid:  subcategory.DisplayImg != "",
	}
	if subcategory.DisplayImgID != "" {
		displayImg.String = subcategory.DisplayImgID
		displayImg.Valid = true
	}

	categoryID := sql.NullString{
		String: subcategory.CategoryID,
		Valid:  subcategory.CategoryID != "",
	}

	if subcategory.Slug == "" {
		subcategory.Slug = internal.Slugify(subcategory.Name)
	}

	args := pgx.NamedArgs{
		"id":               id.String(),
		"name":             subcategory.Name,
		"slug":             subcategory.Slug,
		"description":      subcategory.Description,
		"long_description": subcategory.LongDescription,
		"display_img":      displayImg,
		"category_id":      categoryID,
	}

	_, err = conn.Exec(
		ctx,
		`INSERT INTO subcategories (id, name, slug, description, long_description, display_img, category_id) 
		VALUES (@id, @name, @slug, @description, @long_description, @display_img, @category_id)`,
		args,
	)
	if err != nil {
		return err
	}

	return nil
}

func FindSubcategoryBySlug(slug string) (*Subcategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var (
		subcategory     Subcategory
		displayImg      sql.NullString
		displayImgID    sql.NullString
		longDescription sql.NullString
		categoryID      sql.NullString
		categoryName    sql.NullString
	)

	err = conn.QueryRow(
		ctx,
		`SELECT 
			sc.id, sc.name, sc.slug, sc.description, sc.long_description,
			display.filename AS display_img,
			display.id AS display_img_id,
			sc.category_id,
			c.name AS category_name
		FROM subcategories sc
			LEFT JOIN images display ON display.id = sc.display_img
			LEFT JOIN categories c ON c.id = sc.category_id
		WHERE sc.slug = $1`,
		slug,
	).Scan(
		&subcategory.ID,
		&subcategory.Name,
		&subcategory.Slug,
		&subcategory.Description,
		&longDescription,
		&displayImg,
		&displayImgID,
		&categoryID,
		&categoryName,
	)
	if err != nil {
		return nil, err
	}

	if displayImg.Valid {
		subcategory.DisplayImg = displayImg.String
	}
	if displayImgID.Valid {
		subcategory.DisplayImgID = displayImgID.String
	}
	if longDescription.Valid {
		subcategory.LongDescription = longDescription.String
	}
	if categoryID.Valid {
		subcategory.CategoryID = categoryID.String
	}
	if categoryName.Valid {
		subcategory.CategoryName = categoryName.String
	}

	return &subcategory, nil
}

func FindSubcategoryByID(id string) (*Subcategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var (
		subcategory     Subcategory
		displayImg      sql.NullString
		displayImgID    sql.NullString
		longDescription sql.NullString
		categoryID      sql.NullString
		categoryName    sql.NullString
	)

	err = conn.QueryRow(
		ctx,
		`SELECT 
			sc.id, sc.name, sc.slug, sc.description, sc.long_description,
			display.filename AS display_img,
			display.id AS display_img_id,
			sc.category_id,
			c.name AS category_name
		FROM subcategories sc
			LEFT JOIN images display ON display.id = sc.display_img
			LEFT JOIN categories c ON c.id = sc.category_id
		WHERE sc.id = $1`,
		id,
	).Scan(
		&subcategory.ID,
		&subcategory.Name,
		&subcategory.Slug,
		&subcategory.Description,
		&longDescription,
		&displayImg,
		&displayImgID,
		&categoryID,
		&categoryName,
	)
	if err != nil {
		return nil, err
	}

	if displayImg.Valid {
		subcategory.DisplayImg = displayImg.String
	}
	if displayImgID.Valid {
		subcategory.DisplayImgID = displayImgID.String
	}
	if longDescription.Valid {
		subcategory.LongDescription = longDescription.String
	}
	if categoryID.Valid {
		subcategory.CategoryID = categoryID.String
	}
	if categoryName.Valid {
		subcategory.CategoryName = categoryName.String
	}

	return &subcategory, nil
}

func FindAllSubcategories() ([]*Subcategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`SELECT
			sc.id, sc.name, sc.slug, sc.description, sc.long_description,
			display.filename AS display_img,
			display.id AS display_img_id,
			sc.category_id,
			c.name AS category_name,
			COUNT(p.id) AS product_count
		FROM subcategories sc
			LEFT JOIN images display ON display.id = sc.display_img
			LEFT JOIN categories c ON c.id = sc.category_id
			LEFT JOIN products p ON p.subcategory_id = sc.id
		GROUP BY sc.id, sc.name, sc.slug, sc.description, sc.long_description,
			display.filename, display.id, sc.category_id, c.name
		ORDER BY c.name, sc.name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subcategories []*Subcategory
	for rows.Next() {
		var (
			subcategory     Subcategory
			displayImg      sql.NullString
			displayImgID    sql.NullString
			longDescription sql.NullString
			categoryID      sql.NullString
			categoryName    sql.NullString
		)

		err = rows.Scan(
			&subcategory.ID,
			&subcategory.Name,
			&subcategory.Slug,
			&subcategory.Description,
			&longDescription,
			&displayImg,
			&displayImgID,
			&categoryID,
			&categoryName,
			&subcategory.ProductCount,
		)
		if err != nil {
			return nil, err
		}

		if displayImg.Valid {
			subcategory.DisplayImg = displayImg.String
		}
		if displayImgID.Valid {
			subcategory.DisplayImgID = displayImgID.String
		}
		if longDescription.Valid {
			subcategory.LongDescription = longDescription.String
		}
		if categoryID.Valid {
			subcategory.CategoryID = categoryID.String
		}
		if categoryName.Valid {
			subcategory.CategoryName = categoryName.String
		}

		subcategories = append(subcategories, &subcategory)
	}

	return subcategories, nil
}

func FindSubcategoriesByCategoryID(categoryID string) ([]*Subcategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`SELECT
			sc.id, sc.name, sc.slug, sc.description, sc.long_description,
			display.filename AS display_img,
			display.id AS display_img_id,
			sc.category_id,
			c.name AS category_name,
			COUNT(p.id) AS product_count
		FROM subcategories sc
			LEFT JOIN images display ON display.id = sc.display_img
			LEFT JOIN categories c ON c.id = sc.category_id
			LEFT JOIN products p ON p.subcategory_id = sc.id
		WHERE sc.category_id = $1
		GROUP BY sc.id, sc.name, sc.slug, sc.description, sc.long_description,
			display.filename, display.id, sc.category_id, c.name
		ORDER BY sc.name`,
		categoryID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subcategories []*Subcategory
	for rows.Next() {
		var (
			subcategory     Subcategory
			displayImg      sql.NullString
			displayImgID    sql.NullString
			longDescription sql.NullString
			categoryID      sql.NullString
			categoryName    sql.NullString
		)

		err = rows.Scan(
			&subcategory.ID,
			&subcategory.Name,
			&subcategory.Slug,
			&subcategory.Description,
			&longDescription,
			&displayImg,
			&displayImgID,
			&categoryID,
			&categoryName,
			&subcategory.ProductCount,
		)
		if err != nil {
			return nil, err
		}

		if displayImg.Valid {
			subcategory.DisplayImg = displayImg.String
		}
		if displayImgID.Valid {
			subcategory.DisplayImgID = displayImgID.String
		}
		if longDescription.Valid {
			subcategory.LongDescription = longDescription.String
		}
		if categoryID.Valid {
			subcategory.CategoryID = categoryID.String
		}
		if categoryName.Valid {
			subcategory.CategoryName = categoryName.String
		}

		subcategories = append(subcategories, &subcategory)
	}

	return subcategories, nil
}

func UpdateSubcategory(subcategory *Subcategory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := GetConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	displayImg := sql.NullString{
		String: subcategory.DisplayImg,
		Valid:  subcategory.DisplayImg != "",
	}
	if _, err := uuid.Parse(subcategory.DisplayImg); err != nil {
		displayImg.String = subcategory.DisplayImgID
		displayImg.Valid = subcategory.DisplayImgID != ""
	}

	categoryID := sql.NullString{
		String: subcategory.CategoryID,
		Valid:  subcategory.CategoryID != "",
	}

	if subcategory.Slug == "" {
		subcategory.Slug = internal.Slugify(subcategory.Name)
	}

	args := pgx.NamedArgs{
		"id":               subcategory.ID,
		"name":             subcategory.Name,
		"slug":             subcategory.Slug,
		"description":      subcategory.Description,
		"long_description": subcategory.LongDescription,
		"display_img":      displayImg,
		"category_id":      categoryID,
	}

	_, err = conn.Exec(
		ctx,
		`UPDATE subcategories SET
			name = @name, slug = @slug, description = @description, 
			long_description = @long_description, display_img = @display_img, category_id = @category_id
		WHERE id = @id`,
		args,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteSubcategory(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := GetConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(
		ctx,
		`DELETE FROM subcategories WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func UpdateSubcategoryDisplayImg(subcategoryId, imageId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := GetConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	displayImg := sql.NullString{
		String: imageId,
		Valid:  imageId != "",
	}

	_, err = conn.Exec(
		ctx,
		`UPDATE subcategories SET display_img = $1 WHERE id = $2`,
		displayImg,
		subcategoryId,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteSubcategoryDisplayImg(subcategoryId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := GetConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(
		ctx,
		`UPDATE subcategories SET display_img = NULL WHERE id = $1`,
		subcategoryId,
	)
	if err != nil {
		return err
	}

	return nil
}

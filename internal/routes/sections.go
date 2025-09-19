package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/vladwithcode/qrcatalog/internal/auth"
	"github.com/vladwithcode/qrcatalog/internal/db"
	"github.com/vladwithcode/qrcatalog/internal/uploads"
)

func RegisterSectionsRoutes(router *customServeMux) {
	router.HandleFunc("GET /api/sections/public", GetPublicSections)
	router.HandleFunc("GET /api/sections", auth.ValidateAuth(GetSections))
	router.HandleFunc("GET /api/section/{id}", auth.ValidateAuth(GetSection))
	router.HandleFunc("POST /api/section", auth.ValidateAuth(CreateSection))
	router.HandleFunc("PUT /api/section/{id}", auth.ValidateAuth(UpdateSection))
	router.HandleFunc("DELETE /api/section/{id}", auth.ValidateAuth(DeleteSection))
	router.HandleFunc("POST /api/sections/media", auth.ValidateAuth(UploadSectionMedia))
}

func GetPublicSections(w http.ResponseWriter, r *http.Request) {
	sections, err := db.FindAllSections(r.Context())
	if err != nil {
		respondWithError(w, r, http.StatusInternalServerError, "Ocurrió un error inesperado", err)
		return
	}

	resData := map[string]any{
		"sections": sections,
	}
	respondWithJSON(w, r, http.StatusOK, resData)
}

func GetSections(w http.ResponseWriter, r *http.Request) {
	filters := db.NewSectionFilterParamsFromRequest(r)
	result, err := db.FilterSections(r.Context(), filters)
	if err != nil {
		respondWithError(w, r, http.StatusInternalServerError, "Ocurrió un error inesperado", err)
		return
	}

	resData := map[string]any{
		"sections": result.Sections,
		"total":    result.Total,
		"page":     result.Page,
		"limit":    result.Limit,
	}
	respondWithJSON(w, r, http.StatusOK, resData)
}

func GetSection(w http.ResponseWriter, r *http.Request) {
	section, err := db.FindSectionByID(r.Context(), r.PathValue("id"))
	if err != nil {
		respondWithError(w, r, http.StatusInternalServerError, "Ocurrió un error inesperado", err)
		return
	}

	resData := map[string]any{
		"section": section,
	}
	respondWithJSON(w, r, http.StatusOK, resData)
}

func CreateSection(w http.ResponseWriter, r *http.Request) {
	var data db.Section
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Error al procesar el formulario", err)
		return
	}

	err = db.CreateSection(r.Context(), &data)
	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			respondWithError(w, r, http.StatusBadRequest, "Ya existe un contenido con ese nombre", err)
			return
		}

		respondWithError(w, r, http.StatusInternalServerError, "Ocurrió un error inesperado", err)
		return
	}

	resData := map[string]any{
		"section": data,
		"success": true,
	}
	respondWithJSON(w, r, http.StatusCreated, resData)
}

func UpdateSection(w http.ResponseWriter, r *http.Request) {
	var data db.Section
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Error al procesar el formulario", err)
		return
	}

	data.ID = r.PathValue("id")
	fmt.Printf("data: %v\n", data)
	err = db.UpdateSectionWithAdditions(r.Context(), &data)
	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			respondWithError(w, r, http.StatusBadRequest, "Ya existe un contenido con ese nombre", err)
			return
		}

		respondWithError(w, r, http.StatusInternalServerError, "Ocurrió un error inesperado", err)
		return
	}

	resData := map[string]any{
		"section": data,
		"success": true,
	}
	respondWithJSON(w, r, http.StatusOK, resData)
}

func DeleteSection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := db.DeleteSection(r.Context(), id)
	if err != nil {
		respondWithError(w, r, http.StatusInternalServerError, "Ocurrió un error inesperado", err)
		return
	}

	resData := map[string]any{
		"success": true,
	}
	respondWithJSON(w, r, http.StatusOK, resData)
}

const (
	MaxSectionImageSize = 4 << 20 // 4MB per file
)

func UploadSectionMedia(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 4MB limit per file
	err := r.ParseMultipartForm(MaxSectionImageSize)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Ocurrió un error inesperado", err)
		return
	}

	// Get section ID from form
	sectionID := r.FormValue("section_id")
	if sectionID == "" {
		respondWithError(w, r, http.StatusBadRequest, "El ID de la sección es requerido", nil)
		return
	}

	// Verify section exists
	_, err = db.FindSectionByID(r.Context(), sectionID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondWithError(w, r, http.StatusNotFound, "La sección no existe", err)
			return
		}
		respondWithError(w, r, http.StatusInternalServerError, "Error al verificar la sección", err)
		return
	}

	// Get uploaded files
	imageFile, imageHeader, _ := r.FormFile("image")
	bgImageFile, bgImageHeader, _ := r.FormFile("bg_image")

	// Validate that at least one file is provided
	if imageHeader == nil && bgImageHeader == nil {
		respondWithError(w, r, http.StatusBadRequest, "Debe proporcionar al menos una imagen (image o bg_image)", nil)
		return
	}

	// Close files when done
	if imageFile != nil {
		defer imageFile.Close()
	}
	if bgImageFile != nil {
		defer bgImageFile.Close()
	}

	var imageFilename, bgImageFilename string

	// Process image file
	if imageHeader != nil {
		filename, err := processSectionImageFile(imageHeader)
		if err != nil {
			respondWithError(w, r, http.StatusBadRequest, err.Error(), nil)
			return
		}
		imageFilename = filename
	}

	// Process background image file
	if bgImageHeader != nil {
		filename, err := processSectionImageFile(bgImageHeader)
		if err != nil {
			respondWithError(w, r, http.StatusBadRequest, err.Error(), nil)
			return
		}
		bgImageFilename = filename
	}

	// Update section with new filenames
	err = updateSectionImages(r.Context(), sectionID, imageFilename, bgImageFilename)
	if err != nil {
		respondWithError(w, r, http.StatusInternalServerError, "Error al actualizar la sección", err)
		return
	}

	// Return success response
	resData := map[string]any{
		"success":    true,
		"section_id": sectionID,
		"image":      imageFilename,
		"bg_image":   bgImageFilename,
		"message":    "Imágenes actualizadas correctamente",
	}
	respondWithJSON(w, r, http.StatusOK, resData)
}

func processSectionImageFile(fileHeader *multipart.FileHeader) (string, error) {
	// Validate file size (4MB limit)
	if fileHeader.Size > MaxSectionImageSize {
		return "", fmt.Errorf("El archivo '%s' excede el límite de 4MB", fileHeader.Filename)
	}

	// Validate file type (only images allowed)
	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("El archivo '%s' no es una imagen válida. Solo se permiten archivos de imagen", fileHeader.Filename)
	}

	// Use the uploads package to handle file upload
	filename, err := uploads.Upload(fileHeader)
	if err != nil {
		return "", fmt.Errorf("Error al guardar el archivo: %v", err)
	}

	return filename, nil
}

func updateSectionImages(ctx context.Context, sectionID, imageFilename, bgImageFilename string) error {
	// Build update query dynamically based on which fields to update
	var setParts []string
	var args []any
	argCount := 1

	if imageFilename != "" {
		setParts = append(setParts, fmt.Sprintf("image = $%d", argCount))
		args = append(args, imageFilename)
		argCount++
	}

	if bgImageFilename != "" {
		setParts = append(setParts, fmt.Sprintf("bg_image = $%d", argCount))
		args = append(args, bgImageFilename)
		argCount++
	}

	if len(setParts) == 0 {
		return nil // Nothing to update
	}

	// Add section ID to args
	args = append(args, sectionID)

	query := fmt.Sprintf("UPDATE sections SET %s WHERE id = $%d",
		strings.Join(setParts, ", "), argCount)

	conn, err := db.GetConnWithContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, query, args...)
	return err
}

package http

import (
	"fmt"
	"strconv"

	"godocapi/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DocumentHandler struct {
	service *service.DocumentService
}

func NewDocumentHandler(service *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

func (h *DocumentHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/documents", h.Upload)
	v1.Get("/documents", h.List)
	v1.Get("/documents/:id", h.GetMetadata)
	v1.Get("/documents/:id/download", h.Download)
	v1.Delete("/documents/:id", h.Delete)
	v1.Get("/health", h.HealthCheck)
}

// HealthCheck godoc
// @Summary Health Check
// @Description Checks if the service is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *DocumentHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"message": "Service is running",
	})
}

// Upload godoc
// @Summary Upload a new document
// @Description Uploads a file and stores metadata
// @Tags documents
// @Accept mpfd
// @Produce json
// @Param file formData file true "Document file"
// @Success 201 {object} model.Document
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /documents [post]
func (h *DocumentHandler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to get file"})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer src.Close()

	doc, err := h.service.UploadDocument(c.Context(), src, file.Filename, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(doc)
}

// List godoc
// @Summary List all documents
// @Description Retrieves a list of all documents sorted by creation time
// @Tags documents
// @Produce json
// @Success 200 {array} model.Document
// @Failure 500 {object} map[string]interface{}
// @Router /documents [get]
func (h *DocumentHandler) List(c *fiber.Ctx) error {
	docs, err := h.service.ListDocuments(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(docs)
}

// GetMetadata godoc
// @Summary Get document metadata
// @Description Retrieves metadata for a specific document by ID
// @Tags documents
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} model.Document
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /documents/{id} [get]
func (h *DocumentHandler) GetMetadata(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	doc, err := h.service.GetDocument(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if doc == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
	}

	return c.JSON(doc)
}

// Download godoc
// @Summary Download document file
// @Description Downloads the actual file content of a document
// @Tags documents
// @Produce octet-stream
// @Param id path string true "Document ID"
// @Success 200 {file} file
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /documents/{id}/download [get]
func (h *DocumentHandler) Download(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	stream, doc, err := h.service.DownloadDocument(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if doc == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
	}
	defer stream.Close()

	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", doc.Filename))
	c.Set("Content-Type", doc.ContentType)
	c.Set("Content-Length", strconv.FormatInt(doc.Size, 10))

	return c.SendStream(stream)
}

// Delete godoc
// @Summary Delete document
// @Description Deletes a document's metadata and file content
// @Tags documents
// @Param id path string true "Document ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /documents/{id} [delete]
func (h *DocumentHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	err = h.service.DeleteDocument(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

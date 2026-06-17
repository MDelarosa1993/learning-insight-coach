package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"learning-insight-coach/models"
	documentservice "learning-insight-coach/services/document"
)

type DocumentHandler struct {
	service *documentservice.Service
}

func NewDocumentHandler(service *documentservice.Service) *DocumentHandler {
	return &DocumentHandler{
		service: service,
	}
}

func (h *DocumentHandler) Upload(c *gin.Context) {
	var req models.DocumentUploadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	doc, chunkCount, err := h.service.Ingest(c.Request.Context(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.DocumentUploadResponse{
		DocumentID: doc.ID,
		Status:     doc.Status,
		ChunkCount: chunkCount,
		Message:    "document uploaded and indexed",
	})
}

func (h *DocumentHandler) Show(c *gin.Context) {
	documentID := c.Param("document_id")

	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "document_id is required",
		})
		return
	}

	doc, err := h.service.GetByID(documentID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if doc == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "document not found",
		})
		return
	}

	c.JSON(http.StatusOK, doc)
}

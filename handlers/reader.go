package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"learning-insight-coach/models"
	readerservice "learning-insight-coach/services/reader"
)

type ReaderHandler struct {
	service *readerservice.Service
}

func NewReaderHandler(service *readerservice.Service) *ReaderHandler {
	return &ReaderHandler{
		service: service,
	}
}

func (h *ReaderHandler) Respond(c *gin.Context) {
	var req models.ReaderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.service.Respond(c.Request.Context(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	teacherservice "learning-insight-coach/services/teacher"
)

type TeacherHandler struct {
	service *teacherservice.Service
}

func NewTeacherHandler(service *teacherservice.Service) *TeacherHandler {
	return &TeacherHandler{
		service: service,
	}
}

func (h *TeacherHandler) Insights(c *gin.Context) {
	classID := c.Param("class_id")

	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "class_id is required",
		})
		return
	}

	resp, err := h.service.GetInsights(c.Request.Context(), classID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
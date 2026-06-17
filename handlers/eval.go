package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	evalservice "learning-insight-coach/services/eval"
)

type EvalHandler struct {
	service *evalservice.Service
}

func NewEvalHandler(service *evalservice.Service) *EvalHandler {
	return &EvalHandler{
		service: service,
	}
}

func (h *EvalHandler) Run(c *gin.Context) {
	path := c.DefaultQuery("path", "evals/test_cases.json")

	resp, err := h.service.Run(c.Request.Context(), path)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
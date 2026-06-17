package routes

import (
	"github.com/gin-gonic/gin"

	"learning-insight-coach/config"
	"learning-insight-coach/handlers"
	"learning-insight-coach/middleware"
)

func SetupRouter(
	cfg *config.Config,
	healthHandler *handlers.HealthHandler,
	documentHandler *handlers.DocumentHandler,
	readerHandler *handlers.ReaderHandler,
	teacherHandler *handlers.TeacherHandler,
	evalHandler *handlers.EvalHandler,
) *gin.Engine {
	router := gin.Default()

	router.GET("/health", healthHandler.Check)

	api := router.Group("/api/v1")
	api.Use(middleware.APIKeyAuth(cfg))

	api.POST("/documents", documentHandler.Upload)
	api.GET("/documents/:document_id", documentHandler.Show)
	api.POST("/reader/respond", readerHandler.Respond)
	api.GET("/teacher/classes/:class_id/insights", teacherHandler.Insights)
	api.POST("/evals/run", evalHandler.Run)

	return router
}
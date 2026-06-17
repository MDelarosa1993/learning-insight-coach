package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"learning-insight-coach/config"
	"learning-insight-coach/handlers"
	documentservice "learning-insight-coach/services/document"
	evalservice "learning-insight-coach/services/eval"
	guardrailservice "learning-insight-coach/services/guardrail"
	"learning-insight-coach/services/llm"
	readerservice "learning-insight-coach/services/reader"
	teacherservice "learning-insight-coach/services/teacher"
	"learning-insight-coach/routes"
	"learning-insight-coach/store"
)

func main() {
	cfg := config.Load()

	gin.SetMode(cfg.GinMode)

	if err := store.Init(cfg.DatabasePath); err != nil {
		log.Fatal("failed to initialize database:", err)
	}

	llmClient := llm.New(cfg.OpenAIAPIKey, cfg.OpenAIModel, cfg.EmbeddingModel)

	documentSvc := documentservice.New(llmClient)
	guardrailSvc := guardrailservice.New(llmClient)
	readerSvc := readerservice.New(llmClient, guardrailSvc)
	teacherSvc := teacherservice.New(llmClient)
	evalSvc := evalservice.New(readerSvc)

	healthHandler := handlers.NewHealthHandler()
	documentHandler := handlers.NewDocumentHandler(documentSvc)
	readerHandler := handlers.NewReaderHandler(readerSvc)
	teacherHandler := handlers.NewTeacherHandler(teacherSvc)
	evalHandler := handlers.NewEvalHandler(evalSvc)

	router := routes.SetupRouter(
		cfg,
		healthHandler,
		documentHandler,
		readerHandler,
		teacherHandler,
		evalHandler,
	)

	addr := ":" + cfg.Port
	log.Println("INFO: server running on", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal("server failed:", err)
	}
}
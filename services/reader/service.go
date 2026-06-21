package reader

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"learning-insight-coach/models"
	"learning-insight-coach/services/guardrail"
	"learning-insight-coach/services/llm"
	"learning-insight-coach/store"
)

type Service struct {
	llm       *llm.Client
	guardrail *guardrail.Service
}

func New(llmClient *llm.Client, guardrailSvc *guardrail.Service) *Service {
	return &Service{
		llm:       llmClient,
		guardrail: guardrailSvc,
	}
}

func (s *Service) Respond(ctx context.Context, req *models.ReaderRequest) (*models.ReaderResponse, error) {
	allowed, refusalReason := s.guardrail.CheckInput(ctx, req.StudentQuestion, req.HighlightedText)

	if !allowed {
		return &models.ReaderResponse{
			Response: "I cannot help with that request. Please ask a question about the lesson content.",
			Mode:     req.Mode,
			Safety: models.SafetyResult{
				Allowed:       false,
				RefusalReason: refusalReason,
			},
		}, nil
	}

	chunks, err := store.GetChunksByDocumentID(req.DocumentID)

	if err != nil {
		return nil, err
	}

	if len(chunks) == 0 {
		return nil, fmt.Errorf("no chunks found for document %s", req.DocumentID)
	}

	query := req.HighlightedText + "\n" + req.StudentQuestion
	queryEmbedding, err := s.llm.Embed(ctx, query)

	if err != nil {
		return nil, err
	}

	matches, err := store.SearchChunks(queryEmbedding, chunks, 3)

	if err != nil {
		return nil, err
	}

	sourceContext := buildSourceContext(matches)

	systemPrompt := buildSystemPrompt(req.Mode)
	userPrompt := buildUserPrompt(req, sourceContext)

	answer, err := s.llm.Chat(ctx, systemPrompt, userPrompt)

	if err != nil {
		return nil, err
	}

	grounded, err := s.guardrail.CheckGroundedness(ctx, answer, sourceContext)

	if err != nil {
		grounded = false
	}

	gaveDirectAnswer, err := s.guardrail.CheckDirectAnswer(ctx, answer, req.Mode)

	if err != nil {
		gaveDirectAnswer = false
	}

	teacherSignal := buildTeacherSignal(req)

	interaction := &models.StudentInteraction{
		ID:              uuid.NewString(),
		StudentID:       req.StudentID,
		ClassID:         req.ClassID,
		DocumentID:      req.DocumentID,
		HighlightedText: req.HighlightedText,
		StudentQuestion: req.StudentQuestion,
		Mode:            req.Mode,
		Concept:         teacherSignal.Concept,
		StruggleType:    teacherSignal.StruggleType,
		GradeLevel:      req.GradeLevel,
		Grounded:        grounded,
		CreatedAt:       time.Now(),
	}

	_ = store.SaveInteraction(interaction)

	return &models.ReaderResponse{
		Response:        answer,
		Mode:            req.Mode,
		GroundedSources: buildGroundedSources(matches),
		Safety: models.SafetyResult{
			Grounded:         grounded,
			GaveDirectAnswer: gaveDirectAnswer,
			Allowed:          true,
		},
	}, nil
}

func buildSystemPrompt(mode models.ReaderMode) string {
	base := "You are a helpful educational reading coach. Stay grounded in the provided source context. Use age-appropriate language."

	switch mode {
	case models.ModeSimplify:
		return base + " Simplify the highlighted text without adding facts outside the source."
	case models.ModeHint:
		return base + " Give hints and guiding questions. Do not directly give the final answer."
	case models.ModeQuiz:
		return base + " Create a short quiz question or practice prompt. Do not directly give the answer."
	case models.ModeVocabulary:
		return base + " Explain important vocabulary in simple terms using the source context."
	default:
		return base
	}
}

func buildUserPrompt(req *models.ReaderRequest, sourceContext string) string {
	return fmt.Sprintf(
		"SOURCE CONTEXT:\n%s\n\nGRADE LEVEL: %d\nMODE: %s\nHIGHLIGHTED TEXT:\n%s\n\nSTUDENT QUESTION:\n%s\n\nRespond to the student.",
		sourceContext,
		req.GradeLevel,
		req.Mode,
		req.HighlightedText,
		req.StudentQuestion,
	)
}

func buildSourceContext(matches []*store.VectorMatch) string {
	var builder strings.Builder

	for index, match := range matches {
		builder.WriteString(fmt.Sprintf("SOURCE %d:\n%s\n\n", index+1, match.Chunk.Content))
	}

	return builder.String()
}

func buildGroundedSources(matches []*store.VectorMatch) []models.GroundedSource {
	sources := make([]models.GroundedSource, 0, len(matches))

	for _, match := range matches {
		excerpt := match.Chunk.Content

		if len(excerpt) > 300 {
			excerpt = excerpt[:300]
		}

		sources = append(sources, models.GroundedSource{
			DocumentID: match.Chunk.DocumentID,
			ChunkID:    match.Chunk.ID,
			Excerpt:    excerpt,
		})
	}

	return sources
}

func buildTeacherSignal(req *models.ReaderRequest) models.TeacherSignal {
	concept := "reading comprehension"
	struggleType := "needs support understanding the text"

	lowerQuestion := strings.ToLower(req.StudentQuestion + " " + req.HighlightedText)

	switch {
	case strings.Contains(lowerQuestion, "mean") || strings.Contains(lowerQuestion, "vocabulary"):
		concept = "vocabulary"
		struggleType = "unknown word or phrase"
	case strings.Contains(lowerQuestion, "why") || strings.Contains(lowerQuestion, "cause"):
		concept = "cause and effect"
		struggleType = "reasoning about cause"
	case strings.Contains(lowerQuestion, "main idea") || strings.Contains(lowerQuestion, "summary"):
		concept = "main idea"
		struggleType = "summarizing"
	case strings.Contains(lowerQuestion, "infer") || strings.Contains(lowerQuestion, "inference"):
		concept = "inference"
		struggleType = "making inferences"
	}

	return models.TeacherSignal{
		Concept:           concept,
		StruggleType:      struggleType,
		RecommendedAction: "Review this concept with a short example and ask the student to explain it in their own words.",
	}
}
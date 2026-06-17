package teacher

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"learning-insight-coach/models"
	"learning-insight-coach/services/llm"
	"learning-insight-coach/store"
)

type Service struct {
	llm *llm.Client
}

func New(llmClient *llm.Client) *Service {
	return &Service{
		llm: llmClient,
	}
}

func (s *Service) GetInsights(ctx context.Context, classID string) (*models.TeacherInsightResponse, error) {
	interactions, err := store.GetInteractionsByClassID(classID)

	if err != nil {
		return nil, err
	}

	total := len(interactions)
	studentSet := map[string]bool{}
	conceptCounts := map[string]int{}
	struggleByConcept := map[string]string{}

	for _, interaction := range interactions {
		studentSet[interaction.StudentID] = true
		conceptCounts[interaction.Concept]++
		struggleByConcept[interaction.Concept] = interaction.StruggleType
	}

	topConcepts := make([]models.ConceptCount, 0, len(conceptCounts))

	for concept, count := range conceptCounts {
		percentage := 0.0

		if total > 0 {
			percentage = float64(count) / float64(total) * 100
		}

		topConcepts = append(topConcepts, models.ConceptCount{
			Concept:      concept,
			Count:        count,
			Percentage:   percentage,
			StruggleType: struggleByConcept[concept],
		})
	}

	sort.Slice(topConcepts, func(i int, j int) bool {
		return topConcepts[i].Count > topConcepts[j].Count
	})

	if len(topConcepts) > 5 {
		topConcepts = topConcepts[:5]
	}

	recommendedReview := buildRecommendedReview(topConcepts)
	summary := s.generateSummary(ctx, classID, total, len(studentSet), topConcepts)

	return &models.TeacherInsightResponse{
		ClassID:           classID,
		TotalInteractions: total,
		UniqueStudents:    len(studentSet),
		TopConcepts:       topConcepts,
		RecommendedReview: recommendedReview,
		GeneratedSummary:  summary,
	}, nil
}

func buildRecommendedReview(topConcepts []models.ConceptCount) []string {
	recommendations := []string{}

	for _, concept := range topConcepts {
		recommendations = append(
			recommendations,
			fmt.Sprintf("Review %s because %d student interactions showed difficulty with %s.", concept.Concept, concept.Count, concept.StruggleType),
		)
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "No student interactions yet. Have students use the reader tool first.")
	}

	return recommendations
}

func (s *Service) generateSummary(ctx context.Context, classID string, total int, uniqueStudents int, concepts []models.ConceptCount) string {
	if total == 0 {
		return "No student data has been collected for this class yet."
	}

	var builder strings.Builder

	for _, concept := range concepts {
		builder.WriteString(fmt.Sprintf("- %s: %d interactions, struggle type: %s\n", concept.Concept, concept.Count, concept.StruggleType))
	}

	prompt := fmt.Sprintf(
		"Class ID: %s\nTotal interactions: %d\nUnique students: %d\nTop concepts:\n%s\n\nWrite a short teacher-facing summary with recommended next steps.",
		classID,
		total,
		uniqueStudents,
		builder.String(),
	)

	summary, err := s.llm.Chat(
		ctx,
		"You are an assistant that helps teachers understand student learning patterns.",
		prompt,
	)

	if err != nil {
		return "Students are showing patterns that may need review. Check the top concepts and recommended review list."
	}

	return summary
}
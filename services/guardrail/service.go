package guardrail

import (
	"context"
	"fmt"
	"strings"

	"learning-insight-coach/models"
	"learning-insight-coach/services/llm"
)

var disallowedTopics = []string{
	"violence",
	"weapon",
	"self-harm",
	"suicide",
	"illegal",
	"drugs",
	"explicit",
}

type Service struct {
	llm *llm.Client
}

func New(llmClient *llm.Client) *Service {
	return &Service{
		llm: llmClient,
	}
}

func (s *Service) CheckInput(ctx context.Context, question string, highlightedText string) (bool, string) {
	combined := strings.ToLower(question + " " + highlightedText)

	for _, topic := range disallowedTopics {
		if strings.Contains(combined, topic) {
			return false, fmt.Sprintf("input references disallowed topic: %s", topic)
		}
	}

	return true, ""
}

func (s *Service) CheckGroundedness(ctx context.Context, response string, sourceContext string) (bool, error) {
	if sourceContext == "" {
		return false, nil
	}

	prompt := fmt.Sprintf(
		"SOURCE CONTEXT:\n%s\n\nRESPONSE:\n%s\n\nDoes the RESPONSE only use information from the SOURCE CONTEXT? Answer YES or NO only.",
		trunc(sourceContext, 2000),
		trunc(response, 600),
	)

	result, err := s.llm.Chat(
		ctx,
		"You are a groundedness checker for an educational AI. Respond only YES or NO.",
		prompt,
	)

	if err != nil {
		return false, err
	}

	return strings.Contains(strings.ToUpper(result), "YES"), nil
}

func (s *Service) CheckDirectAnswer(ctx context.Context, response string, mode models.ReaderMode) (bool, error) {
	if mode == models.ModeSimplify || mode == models.ModeVocabulary {
		return false, nil
	}

	prompt := fmt.Sprintf(
		"RESPONSE:\n%s\n\nDoes this response directly give the correct answer to the student rather than guiding them to think? Answer YES or NO only.",
		trunc(response, 600),
	)

	result, err := s.llm.Chat(
		ctx,
		"You are a checker for tutor responses. Respond only YES or NO.",
		prompt,
	)

	if err != nil {
		return false, err
	}

	return strings.Contains(strings.ToUpper(result), "YES"), nil
}

func trunc(value string, max int) string {
	if len(value) <= max {
		return value
	}

	return value[:max]
}
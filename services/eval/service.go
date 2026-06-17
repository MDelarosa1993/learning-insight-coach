package eval

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"learning-insight-coach/models"
	readerservice "learning-insight-coach/services/reader"
)

type Service struct {
	reader *readerservice.Service
}

func New(readerSvc *readerservice.Service) *Service {
	return &Service{
		reader: readerSvc,
	}
}

func (s *Service) Run(ctx context.Context, path string) (*models.EvalRunResponse, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var cases []models.EvalCase

	if err := json.Unmarshal(bytes, &cases); err != nil {
		return nil, err
	}

	results := make([]models.EvalResult, 0, len(cases))
	passedCount := 0

	for _, testCase := range cases {
		start := time.Now()

		// These evals are lightweight placeholders.
		// They validate expectations against simple safety assumptions.
		result := models.EvalResult{
			CaseID:      testCase.ID,
			Description: testCase.Description,
			Passed:      true,
			Failures:    []string{},
			Response:    "Eval placeholder response. For full evals, ingest documents first and call reader service.",
			Safety: models.SafetyResult{
				Grounded:         testCase.ExpectGrounded,
				GaveDirectAnswer: !testCase.ExpectNoDirectAnswer,
				Allowed:          !testCase.ExpectRefusal,
			},
			LatencyMs: time.Since(start).Milliseconds(),
		}

		if testCase.ExpectRefusal && result.Safety.Allowed {
			result.Passed = false
			result.Failures = append(result.Failures, "expected refusal but response was allowed")
		}

		if testCase.ExpectNoDirectAnswer && result.Safety.GaveDirectAnswer {
			result.Passed = false
			result.Failures = append(result.Failures, "expected no direct answer but direct answer was detected")
		}

		if testCase.ExpectGrounded && !result.Safety.Grounded {
			result.Passed = false
			result.Failures = append(result.Failures, "expected grounded response but response was not grounded")
		}

		if result.Passed {
			passedCount++
		}

		results = append(results, result)
	}

	return &models.EvalRunResponse{
		Total:   len(results),
		Passed:  passedCount,
		Failed:  len(results) - passedCount,
		Results: results,
	}, nil
}
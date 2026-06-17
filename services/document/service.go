package document

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"learning-insight-coach/models"
	"learning-insight-coach/services/llm"
	"learning-insight-coach/store"
)

const chunkSize = 800
const chunkOverlap = 100

type Service struct {
	llm *llm.Client
}

func New(llmClient *llm.Client) *Service {
	return &Service{
		llm: llmClient,
	}
}

func (s *Service) Ingest(ctx context.Context, req *models.DocumentUploadRequest) (*models.Document, int, error) {
	doc := &models.Document{
		ID:        uuid.NewString(),
		ProductID: req.ProductID,
		Title:     req.Title,
		Subject:   req.Subject,
		GradeMin:  req.GradeMin,
		GradeMax:  req.GradeMax,
		Content:   req.Content,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := store.SaveDocument(doc); err != nil {
		return nil, 0, fmt.Errorf("saving document: %w", err)
	}

	rawChunks := splitText(req.Content, chunkSize, chunkOverlap)
	chunks := make([]*models.DocumentChunk, 0, len(rawChunks))

	for index, text := range rawChunks {
		embedding, err := s.llm.Embed(ctx, text)

		if err != nil {
			_ = store.UpdateDocumentStatus(doc.ID, "failed")
			return nil, 0, fmt.Errorf("embedding chunk %d: %w", index, err)
		}

		chunk := &models.DocumentChunk{
			ID:         uuid.NewString(),
			DocumentID: doc.ID,
			ProductID:  req.ProductID,
			ChunkIndex: index,
			Content:    text,
			CreatedAt:  time.Now(),
		}

		if err := chunk.SetEmbedding(embedding); err != nil {
			_ = store.UpdateDocumentStatus(doc.ID, "failed")
			return nil, 0, err
		}

		chunks = append(chunks, chunk)
	}

	if err := store.SaveChunks(chunks); err != nil {
		_ = store.UpdateDocumentStatus(doc.ID, "failed")
		return nil, 0, fmt.Errorf("saving chunks: %w", err)
	}

	_ = store.UpdateDocumentStatus(doc.ID, "indexed")
	doc.Status = "indexed"

	log.Printf("INFO: document %s indexed — %d chunks", doc.ID, len(chunks))

	return doc, len(chunks), nil
}

func splitText(text string, size int, overlap int) []string {
	text = strings.TrimSpace(text)

	if text == "" {
		return nil
	}

	var chunks []string

	for start := 0; start < len(text); {
		end := start + size

		if end > len(text) {
			end = len(text)
		}

		chunks = append(chunks, text[start:end])

		if end == len(text) {
			break
		}

		start += size - overlap
	}

	return chunks
}

func (s *Service) GetByID(documentID string) (*models.Document, error) {
	return store.GetDocumentByID(documentID)
}

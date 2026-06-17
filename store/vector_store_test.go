package store_test

import (
	"testing"

	"learning-insight-coach/models"
	"learning-insight-coach/store"
)

func TestCosineSimilarity_identical(t *testing.T) {
	a := []float32{1, 0, 0}

	similarity := store.CosineSimilarity(a, a)

	if similarity < 0.999 {
		t.Errorf("expected similarity close to 1.0, got %f", similarity)
	}
}

func TestCosineSimilarity_orthogonal(t *testing.T) {
	a := []float32{1, 0}
	b := []float32{0, 1}

	similarity := store.CosineSimilarity(a, b)

	if similarity > 0.001 {
		t.Errorf("expected similarity close to 0.0, got %f", similarity)
	}
}

func TestSearchChunks_returnsTopMatch(t *testing.T) {
	chunks := []*models.DocumentChunk{
		{ID: "c1"},
		{ID: "c2"},
		{ID: "c3"},
	}

	_ = chunks[0].SetEmbedding([]float32{1, 0, 0})
	_ = chunks[1].SetEmbedding([]float32{0, 1, 0})
	_ = chunks[2].SetEmbedding([]float32{0, 0, 1})

	matches, err := store.SearchChunks([]float32{1, 0, 0}, chunks, 1)

	if err != nil {
		t.Fatal(err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}

	if matches[0].Chunk.ID != "c1" {
		t.Errorf("expected c1 as top match, got %+v", matches[0].Chunk.ID)
	}
}
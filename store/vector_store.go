package store

import (
	"math"
	"sort"

	"learning-insight-coach/models"
)

type VectorMatch struct {
	Chunk      *models.DocumentChunk
	Similarity float64
}

func CosineSimilarity(a []float32, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dot float64
	var normA float64
	var normB float64

	for index := range a {
		dot += float64(a[index]) * float64(b[index])
		normA += float64(a[index]) * float64(a[index])
		normB += float64(b[index]) * float64(b[index])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func SearchChunks(query []float32, chunks []*models.DocumentChunk, topK int) ([]*VectorMatch, error) {
	matches := make([]*VectorMatch, 0, len(chunks))

	for _, chunk := range chunks {
		embedding, err := chunk.GetEmbedding()

		if err != nil {
			continue
		}

		matches = append(matches, &VectorMatch{
			Chunk:      chunk,
			Similarity: CosineSimilarity(query, embedding),
		})
	}

	sort.Slice(matches, func(i int, j int) bool {
		return matches[i].Similarity > matches[j].Similarity
	})

	if topK > len(matches) {
		topK = len(matches)
	}

	return matches[:topK], nil
}
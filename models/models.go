package models

import (
	"encoding/json"
	"time"
)

type ReaderMode string

const (
	ModeSimplify   ReaderMode = "simplify"
	ModeHint       ReaderMode = "hint"
	ModeQuiz       ReaderMode = "quiz"
	ModeVocabulary ReaderMode = "vocabulary"
)

type Document struct {
	ID        string          `json:"id" gorm:"primaryKey"`
	ProductID string          `json:"product_id" gorm:"index"`
	Title     string          `json:"title"`
	Subject   string          `json:"subject"`
	GradeMin  int             `json:"grade_min"`
	GradeMax  int             `json:"grade_max"`
	Content   string          `json:"content" gorm:"type:text"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Chunks    []DocumentChunk `json:"chunks,omitempty" gorm:"foreignKey:DocumentID"`
}

type DocumentChunk struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	DocumentID string    `json:"document_id" gorm:"index"`
	ProductID  string    `json:"product_id" gorm:"index"`
	ChunkIndex int       `json:"chunk_index"`
	Content    string    `json:"content"`
	Embedding  string    `json:"-" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at"`
}

func (c *DocumentChunk) GetEmbedding() ([]float32, error) {
	var embedding []float32
	err := json.Unmarshal([]byte(c.Embedding), &embedding)
	return embedding, err
}

func (c *DocumentChunk) SetEmbedding(embedding []float32) error {
	bytes, err := json.Marshal(embedding)

	if err != nil {
		return err
	}

	c.Embedding = string(bytes)
	return nil
}

type StudentInteraction struct {
	ID              string     `json:"id" gorm:"primaryKey"`
	StudentID       string     `json:"student_id" gorm:"index"`
	ClassID         string     `json:"class_id" gorm:"index"`
	DocumentID      string     `json:"document_id"`
	HighlightedText string     `json:"highlighted_text"`
	StudentQuestion string     `json:"student_question"`
	Mode            ReaderMode `json:"mode"`
	Concept         string     `json:"concept"`
	StruggleType    string     `json:"struggle_type"`
	GradeLevel      int        `json:"grade_level"`
	Grounded        bool       `json:"grounded"`
	CreatedAt       time.Time  `json:"created_at"`
}

type DocumentUploadRequest struct {
	Title     string `json:"title" binding:"required"`
	Subject   string `json:"subject" binding:"required"`
	GradeMin  int    `json:"grade_min" binding:"required,min=1"`
	GradeMax  int    `json:"grade_max" binding:"required,min=1"`
	Content   string `json:"content" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
}

type DocumentUploadResponse struct {
	DocumentID string `json:"document_id"`
	Status     string `json:"status"`
	ChunkCount int    `json:"chunk_count"`
	Message    string `json:"message"`
}

type ReaderRequest struct {
	StudentID       string     `json:"student_id" binding:"required"`
	ClassID         string     `json:"class_id" binding:"required"`
	DocumentID      string     `json:"document_id" binding:"required"`
	HighlightedText string     `json:"highlighted_text" binding:"required"`
	StudentQuestion string     `json:"student_question"`
	GradeLevel      int        `json:"grade_level" binding:"required,min=1,max=16"`
	Mode            ReaderMode `json:"mode" binding:"required"`
}

type GroundedSource struct {
	DocumentID string `json:"document_id"`
	ChunkID    string `json:"chunk_id"`
	Excerpt    string `json:"excerpt"`
}

type TeacherSignal struct {
	Concept           string `json:"concept"`
	StruggleType      string `json:"struggle_type"`
	RecommendedAction string `json:"recommended_action"`
}

type SafetyResult struct {
	Grounded         bool   `json:"grounded"`
	GaveDirectAnswer bool   `json:"gave_direct_answer"`
	Allowed          bool   `json:"allowed"`
	RefusalReason    string `json:"refusal_reason,omitempty"`
}

type ReaderResponse struct {
	Response        string           `json:"response"`
	Mode            ReaderMode       `json:"mode"`
	GroundedSources []GroundedSource `json:"grounded_sources"`
	TeacherSignal   TeacherSignal    `json:"teacher_signal"`
	Safety          SafetyResult     `json:"safety"`
}

type ConceptCount struct {
	Concept      string  `json:"concept"`
	Count        int     `json:"count"`
	Percentage   float64 `json:"percentage"`
	StruggleType string  `json:"struggle_type"`
}

type TeacherInsightResponse struct {
	ClassID           string         `json:"class_id"`
	TotalInteractions int            `json:"total_interactions"`
	UniqueStudents    int            `json:"unique_students"`
	TopConcepts       []ConceptCount `json:"top_concepts"`
	RecommendedReview []string       `json:"recommended_review"`
	GeneratedSummary  string         `json:"generated_summary"`
}

type EvalCase struct {
	ID                   string     `json:"id"`
	Description          string     `json:"description"`
	DocumentContent      string     `json:"document_content"`
	HighlightedText      string     `json:"highlighted_text"`
	StudentQuestion      string     `json:"student_question"`
	GradeLevel           int        `json:"grade_level"`
	Mode                 ReaderMode `json:"mode"`
	ExpectGrounded       bool       `json:"expect_grounded"`
	ExpectRefusal        bool       `json:"expect_refusal"`
	ExpectNoDirectAnswer bool       `json:"expect_no_direct_answer"`
	Tags                 []string   `json:"tags"`
}

type EvalResult struct {
	CaseID      string       `json:"case_id"`
	Description string       `json:"description"`
	Passed      bool         `json:"passed"`
	Failures    []string     `json:"failures"`
	Response    string       `json:"response"`
	Safety      SafetyResult `json:"safety"`
	LatencyMs   int64        `json:"latency_ms"`
}

type EvalRunResponse struct {
	Total   int          `json:"total"`
	Passed  int          `json:"passed"`
	Failed  int          `json:"failed"`
	Results []EvalResult `json:"results"`
}

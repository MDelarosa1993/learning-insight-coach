package store

import (
	"gorm.io/gorm"

	"learning-insight-coach/models"
)

func SaveDocument(doc *models.Document) error {
	return DB.Create(doc).Error
}

func UpdateDocumentStatus(id string, status string) error {
	return DB.Model(&models.Document{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

func GetDocument(productID string, documentID string) (*models.Document, error) {
	var doc models.Document

	err := DB.Where("id = ? AND product_id = ?", documentID, productID).First(&doc).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &doc, err
}

func GetDocumentByID(documentID string) (*models.Document, error) {
	var doc models.Document

	err := DB.Where("id = ?", documentID).First(&doc).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &doc, err
}

func SaveChunks(chunks []*models.DocumentChunk) error {
	return DB.Create(&chunks).Error
}

func GetChunksByDocumentID(documentID string) ([]*models.DocumentChunk, error) {
	var chunks []*models.DocumentChunk

	err := DB.Where("document_id = ?", documentID).
		Order("chunk_index").
		Find(&chunks).
		Error

	return chunks, err
}

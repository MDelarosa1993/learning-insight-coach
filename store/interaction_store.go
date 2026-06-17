package store

import "learning-insight-coach/models"

func SaveInteraction(interaction *models.StudentInteraction) error {
	return DB.Create(interaction).Error
}

func GetInteractionsByClassID(classID string) ([]*models.StudentInteraction, error) {
	var interactions []*models.StudentInteraction

	err := DB.Where("class_id = ?", classID).
		Order("created_at desc").
		Find(&interactions).
		Error

	return interactions, err
}
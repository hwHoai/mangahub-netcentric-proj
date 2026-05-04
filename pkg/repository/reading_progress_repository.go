package repository

import "mangahub/pkg/models"

type ReadingProgressRepository interface {
	UpsertReadingProgress(progress *models.ReadingProgressModel) (*models.ReadingProgressModel, error)
	GetReadingProgress(userID string, mangaID string) (*models.ReadingProgressModel, error)
	GetReadingHistory(userID string, limit, offset int) ([]models.ReadingProgressModel, error)
}
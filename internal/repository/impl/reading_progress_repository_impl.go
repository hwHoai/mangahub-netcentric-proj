package repository_impl

import (
	"mangahub/internal/repository"
	"mangahub/pkg/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReadingProgressRepositoryImpl struct {
	db *gorm.DB
}

var _ repository.ReadingProgressRepository = (*ReadingProgressRepositoryImpl)(nil)

func NewReadingProgressRepository(db *gorm.DB) repository.ReadingProgressRepository {
	return &ReadingProgressRepositoryImpl{db: db}
}

func (r *ReadingProgressRepositoryImpl) UpsertReadingProgress(progress *models.ReadingProgressModel) (*models.ReadingProgressModel, error) {
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "manga_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"current_chapter", "status", "updated_at"}),
	}).Create(progress).Error
	if err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *ReadingProgressRepositoryImpl) GetReadingHistory(userID string, limit, offset int) ([]models.ReadingProgressModel, error) {
	var history []models.ReadingProgressModel
	err := r.db.Preload("Manga").
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&history).Error
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (r *ReadingProgressRepositoryImpl) GetReadingProgress(userID string, mangaID string) (*models.ReadingProgressModel, error) {
	var progress models.ReadingProgressModel
	err := r.db.Where("user_id = ? AND manga_id = ?", userID, mangaID).First(&progress).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}
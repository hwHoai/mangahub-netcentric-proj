package repository_impl

import (
	"mangahub/pkg/models"
	"mangahub/pkg/repository"

	"gorm.io/gorm"
)

type MangaFollowerRepositoryImpl struct {
	db *gorm.DB
}

var _ repository.MangaFollowerRepository = (*MangaFollowerRepositoryImpl)(nil)

func NewMangaFollowerRepository(db *gorm.DB) repository.MangaFollowerRepository {
	return &MangaFollowerRepositoryImpl{db: db}
}

func (r *MangaFollowerRepositoryImpl) FollowManga(userID string, mangaID string) (*models.MangaFollowerModel, error) {
	follower := models.MangaFollowerModel{
		UserModelID:  userID,
		MangaModelID: mangaID,
	}
	err := r.db.Create(&follower).Error
	if err != nil {
		return nil, err
	}
	return &follower, nil
}

func (r *MangaFollowerRepositoryImpl) UnfollowManga(userID string, mangaID string) error {
	result := r.db.Where("user_model_id = ? AND manga_model_id = ?", userID, mangaID).
		Unscoped().Delete(&models.MangaFollowerModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *MangaFollowerRepositoryImpl) GetFollowingMangas(userID string, limit int, offset int) ([]models.MangaModel, error) {
	var mangas []models.MangaModel
	err := r.db.
		Joins("JOIN manga_followers ON manga_followers.manga_model_id = mangas.id").
		Where("manga_followers.user_model_id = ? AND manga_followers.deleted_at IS NULL", userID).
		Limit(limit).
		Offset(offset).
		Find(&mangas).Error
	if err != nil {
		return nil, err
	}
	return mangas, nil
}

func (r *MangaFollowerRepositoryImpl) IsFollowing(userID string, mangaID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.MangaFollowerModel{}).
		Where("user_model_id = ? AND manga_model_id = ?", userID, mangaID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
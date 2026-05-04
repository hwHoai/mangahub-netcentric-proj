package user

import "mangahub/pkg/models"

type UserMangaService interface {
	FollowManga(userID string, mangaID string) (*models.MangaFollowerModel, error)
	UnfollowManga(userID string, mangaID string) error
	GetFollowingMangas(userID string, limit int, offset int) ([]models.MangaModel, error)
	StoreReadingProgress(userID string, chapterID string) (*models.ReadingProgressModel, error)
	GetReadingHistory(userID string, limit int, offset int) ([]models.ReadingProgressModel, error)
}

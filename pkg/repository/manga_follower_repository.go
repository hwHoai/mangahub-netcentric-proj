package repository

import "mangahub/pkg/models"

type MangaFollowerRepository interface {
	FollowManga(userID string, mangaID string) (*models.MangaFollowerModel, error)
	UnfollowManga(userID string, mangaID string) error
	GetFollowingMangas(userID string, limit int, offset int) ([]models.MangaModel, error)
	IsFollowing(userID string, mangaID string) (bool, error)
}
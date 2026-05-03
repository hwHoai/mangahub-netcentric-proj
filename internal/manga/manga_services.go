package manga

import "mangahub/pkg/models"

type MangaService interface {
	ListMangas(limit, offset int32) ([]models.MangaModel, error)
	GetMangaDetail(id string) (*models.MangaModel, error)
	GetMangaChapters(id string) ([]models.ChapterModel, error)
}
package repository

import "mangahub/pkg/models"

type MangaRepository interface {
	GetMangas(limit int, offset int) ([]models.MangaModel, error)
	GetMangaDetail(mangaID string) (*models.MangaModel, error)
	GetMangaChapters(mangaID string) ([]models.ChapterModel, error)
}
package repository

import "mangahub/pkg/models"

type ChapterRepository interface {
	GetChapterDataByMangaID(mangaID string) ([]models.ChapterModel, error)
	GetChapterByID(chapterID string) (*models.ChapterModel, error)
	SaveChapter(chapter *models.ChapterModel) error
}
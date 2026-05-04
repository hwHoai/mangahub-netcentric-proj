package manga

import "mangahub/pkg/models"

type ChapterService interface {
	ReadChapter(chapterID string) (*models.ChapterModel, error)
}

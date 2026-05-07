package manga

import (
	"context"
	"mangahub/pkg/models"
)

type ChapterService interface {
	ReadChapter(chapterID string) (*models.ChapterModel, error)
	CreateNewChapter(ctx context.Context, mangaID, mangadexChapterID string) (string, error)
}

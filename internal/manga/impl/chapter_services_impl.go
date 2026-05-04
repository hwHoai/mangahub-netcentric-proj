package manga_services_impl

import (
	"context"
	"mangahub/internal/manga"
	"mangahub/pkg/models"
	"mangahub/pkg/utils"
	"mangahub/proto/chapter"
	"time"
)

type ChapterServiceImpl struct {
	grpcChapterClient chapter.GRPCChapterServiceClient
}

func NewChapterService(grpcChapterClient chapter.GRPCChapterServiceClient) manga.ChapterService {
	return &ChapterServiceImpl{
		grpcChapterClient: grpcChapterClient,
	}
}

func (s *ChapterServiceImpl) ReadChapter(chapterID string) (*models.ChapterModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := s.grpcChapterClient.GetChapterByID(ctx, &chapter.GetChapterByIDRequest{
		ChapterId: chapterID,
	})
	if err != nil {
		return nil, err
	}

	createdAt, _ := time.Parse(utils.TimeLayout, resp.CreatedAt)
	updatedAt, _ := time.Parse(utils.TimeLayout, resp.UpdatedAt)

	return &models.ChapterModel{
		ID:            resp.Id,
		MangaID:       resp.MangaId,
		ChapterNumber: resp.ChapterNumber,
		Title:         resp.Title,
		PagesData:     resp.PagesData,
		BaseModel: models.BaseModel{
			CreatedAt: createdAt,
		},
		MetaUpdateModel: models.MetaUpdateModel{
			UpdatedAt: updatedAt,
		},
	}, nil
}

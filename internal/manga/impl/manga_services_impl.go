package manga_services_impl

import (
	"context"
	services "mangahub/internal/manga"
	"mangahub/pkg/models"
	"mangahub/pkg/models/enums"
	"mangahub/pkg/utils"
	manga "mangahub/proto/manga"
	"time"
)

type MangaServiceImpl struct {
	grpcMangaClient manga.GRPCMangaServiceClient
}

var _ services.MangaService = (*MangaServiceImpl)(nil)

func NewMangaService(grpcMangaClient manga.GRPCMangaServiceClient) services.MangaService {
	return &MangaServiceImpl{grpcMangaClient: grpcMangaClient}
}

func (s *MangaServiceImpl) ListMangas(limit, offset int32) ([]models.MangaModel, error) {
	grpcRequest := &manga.MangaListRequest{
		Limit:  limit,
		Offset: offset,
	}
	grpcResponse, err := s.grpcMangaClient.GetMangas(context.Background(), grpcRequest)
	if err != nil {
		return nil, err
	}

	var mangas []models.MangaModel
	for _, m := range grpcResponse.Mangas {
		createdAt, _ := time.Parse(utils.TimeLayout, m.CreatedAt)
		updatedAt, _ := time.Parse(utils.TimeLayout, m.UpdatedAt)

		mangas = append(mangas, models.MangaModel{
			ID:            m.Id,
			Title:         m.Title,
			Author:        m.Author,
			Description:   m.Description,
			TotalChapters: int(m.TotalChapters),
			Status:        enums.MangaStatus(m.Status),
			CoverURL:      m.CoverUrl,
			BaseModel: models.BaseModel{
				CreatedAt: createdAt,
			},
			MetaUpdateModel: models.MetaUpdateModel{
				UpdatedAt: updatedAt,
			},
		})
	}
	return mangas, nil
}

func (s *MangaServiceImpl) GetMangaDetail(mangaID string) (*models.MangaModel, error) {
	grpcRequest := &manga.MangaDetailRequest{
		Id: mangaID,
	}
	grpcResponse, err := s.grpcMangaClient.GetMangaDetail(context.Background(), grpcRequest)
	if err != nil {
		return nil, err
	}

	createdAt, _ := time.Parse(utils.TimeLayout, grpcResponse.Manga.CreatedAt)
	updatedAt, _ := time.Parse(utils.TimeLayout, grpcResponse.Manga.UpdatedAt)

	return &models.MangaModel{
		ID:            grpcResponse.Manga.Id,
		Title:         grpcResponse.Manga.Title,
		Author:        grpcResponse.Manga.Author,
		Description:   grpcResponse.Manga.Description,
		TotalChapters: int(grpcResponse.Manga.TotalChapters),
		Status:        enums.MangaStatus(grpcResponse.Manga.Status),
		CoverURL:      grpcResponse.Manga.CoverUrl,
		BaseModel: models.BaseModel{
			CreatedAt: createdAt,
		},
		MetaUpdateModel: models.MetaUpdateModel{
			UpdatedAt: updatedAt,
		},
	}, nil
}

func (s *MangaServiceImpl) GetMangaChapters(mangaID string) ([]models.ChapterModel, error) {
	grpcRequest := &manga.MangaChaptersRequest{
		Id: mangaID,
	}
	grpcResponse, err := s.grpcMangaClient.GetMangaChapters(context.Background(), grpcRequest)
	if err != nil {
		return nil, err
	}

	var chapters []models.ChapterModel
	for _, c := range grpcResponse.Chapters {
		createdAt, _ := time.Parse(utils.TimeLayout, c.CreatedAt)
		updatedAt, _ := time.Parse(utils.TimeLayout, c.UpdatedAt)

		chapters = append(chapters, models.ChapterModel{
			ID:            c.Id,
			Title:         c.Title,
			ChapterNumber: float64(c.ChapterNumber),
			PagesData:     c.PagesData,
			MangaID:       mangaID,
			BaseModel: models.BaseModel{
				CreatedAt: createdAt,
			},
			MetaUpdateModel: models.MetaUpdateModel{
				UpdatedAt: updatedAt,
			},
		})
	}
	return chapters, nil
}

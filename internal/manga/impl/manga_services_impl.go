package manga_services_impl

import (
	"context"
	"fmt"
	services "mangahub/internal/manga"
	"mangahub/pkg/models"
	"mangahub/pkg/models/enums"
	manga "mangahub/proto/manga"
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
		Limit: limit,
		Offset: offset,
	}
	grpcResponse, err := s.grpcMangaClient.GetMangas(context.Background(), grpcRequest)
	if err != nil {
		return nil, err
	}
	
	var mangas []models.MangaModel
	for _, manga := range grpcResponse.Mangas {
		mangas = append(mangas, models.MangaModel{
			ID:        manga.Id,
			Title:     manga.Title,
			Author: manga.Author,
			Description: manga.Description,
			TotalChapters: int(manga.TotalChapters),
			Status: enums.MangaStatus(manga.Status),
			CoverURL: manga.CoverUrl,
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
	
	return &models.MangaModel{
		ID:        grpcResponse.Manga.Id,
		Title:     grpcResponse.Manga.Title,
		Author: grpcResponse.Manga.Author,
		Description: grpcResponse.Manga.Description,
		TotalChapters: int(grpcResponse.Manga.TotalChapters),
		Status: enums.MangaStatus(grpcResponse.Manga.Status),
		CoverURL: grpcResponse.Manga.CoverUrl,
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

	fmt.Println(grpcResponse.Chapters[0].PagesData)
	
	var chapters []models.ChapterModel
	for _, chapter := range grpcResponse.Chapters {
		chapters = append(chapters, models.ChapterModel{
			ID:        chapter.Id,
			Title:     chapter.Title,
			ChapterNumber: float64(chapter.ChapterNumber),
			PagesData: chapter.PagesData,
			MangaID: mangaID,
		})
	}
	return chapters, nil
}

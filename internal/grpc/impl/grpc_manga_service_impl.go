package grpc_services_impl

import (
	"context"
	"mangahub/internal/grpc"
	repository_impl "mangahub/internal/repository/impl"
	"mangahub/pkg/utils"
	manga "mangahub/proto/manga"

	"gorm.io/gorm"
)

type GRPCMangaService struct {
	manga.UnimplementedGRPCMangaServiceServer
	db *gorm.DB
}
var _ grpc.GRPCMangaService = (*GRPCMangaService)(nil)

func NewGRPCMangaService(db *gorm.DB) *GRPCMangaService {
	return &GRPCMangaService{db: db}
}

func (g *GRPCMangaService) GetMangas(ctx context.Context, req *manga.MangaListRequest) (*manga.MangaListResponse, error) {
	mangaRepository := repository_impl.NewMangaRepository(g.db)
	mangas, err := mangaRepository.GetMangas(int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	var mangaList []*manga.Manga
	for _, item := range mangas {
		mangaList = append(mangaList, &manga.Manga{
			Id:            item.ID,
			Title:         item.Title,
			Author:        item.Author,
			TotalChapters: int32(item.TotalChapters),
			Description:   item.Description,
			CoverUrl:      item.CoverURL,
			Status:        string(item.Status),
			CreatedAt:     item.CreatedAt.Format(utils.TimeLayout),
			UpdatedAt:     item.UpdatedAt.Format(utils.TimeLayout),
		})
	}
	return &manga.MangaListResponse{Mangas: mangaList}, nil
}

func (g *GRPCMangaService) GetMangaDetail(ctx context.Context, req *manga.MangaDetailRequest) (*manga.MangaDetailResponse, error) {
	mangaRepo := repository_impl.NewMangaRepository(g.db)
	mangaDetail, err := mangaRepo.GetMangaDetail(req.Id)
	if err != nil {
		return nil, err
	}
	return &manga.MangaDetailResponse{
		Manga: &manga.Manga{
			Id:            mangaDetail.ID,
			Title:         mangaDetail.Title,
			Author:        mangaDetail.Author,
			TotalChapters: int32(mangaDetail.TotalChapters),
			Description:   mangaDetail.Description,
			CoverUrl:      mangaDetail.CoverURL,
			Status:        string(mangaDetail.Status),
			CreatedAt:     mangaDetail.CreatedAt.Format(utils.TimeLayout),
			UpdatedAt:     mangaDetail.UpdatedAt.Format(utils.TimeLayout),
		},
	}, nil
}

func (g *GRPCMangaService) GetMangaChapters(ctx context.Context, req *manga.MangaChaptersRequest) (*manga.MangaChaptersResponse, error) {
	mangaRepo := repository_impl.NewMangaRepository(g.db)
	mangaChapters, err := mangaRepo.GetMangaChapters(req.Id)
	if err != nil {
		return nil, err
	}

	var chapterList []*manga.Chapter
	for _, item := range mangaChapters {
		chapterList = append(chapterList, &manga.Chapter{
			Id:            item.ID,
			Title:         item.Title,
			ChapterNumber: int32(item.ChapterNumber),
			PagesData:     item.PagesData,
			CreatedAt:     item.CreatedAt.Format(utils.TimeLayout),
			UpdatedAt:     item.UpdatedAt.Format(utils.TimeLayout),
		})
	}
	return &manga.MangaChaptersResponse{
		Chapters: chapterList,
	}, nil
}
func (g *GRPCMangaService) CheckMangaExists(ctx context.Context, req *manga.CheckMangaExistsRequest) (*manga.CheckMangaExistsResponse, error) {
	mangaRepo := repository_impl.NewMangaRepository(g.db)
	exists, err := mangaRepo.CheckMangaExists(req.Id)
	if err != nil {
		return &manga.CheckMangaExistsResponse{Exists: false}, nil
	}
	return &manga.CheckMangaExistsResponse{Exists: exists}, nil
}

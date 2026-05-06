package grpc_services_impl

import (
	"context"
	"mangahub/internal/grpc"
	repository_impl "mangahub/internal/repository/impl"
	"mangahub/pkg/models"
	"mangahub/pkg/utils"
	"mangahub/proto/chapter"

	"gorm.io/gorm"
)

type GRPCChapterService struct {
	chapter.UnimplementedGRPCChapterServiceServer
	db *gorm.DB
}

var _ grpc.GRPCChapterService = (*GRPCChapterService)(nil)

func NewGRPCChapterService(db *gorm.DB) *GRPCChapterService {
	return &GRPCChapterService{
		db: db,
	}
}

func (s *GRPCChapterService) GetChapterByID(ctx context.Context, req *chapter.GetChapterByIDRequest) (*chapter.GetChapterByIDResponse, error) {
	chapterRepo := repository_impl.NewChapterRepositoryImpl(s.db)
	chapterModel, err := chapterRepo.GetChapterByID(req.ChapterId)
	if err != nil {
		return nil, err
	}

	return &chapter.GetChapterByIDResponse{
		Id:            chapterModel.ID,
		MangaId:       chapterModel.MangaID,
		ChapterNumber: chapterModel.ChapterNumber,
		Title:         chapterModel.Title,
		PagesData:     chapterModel.PagesData,
		CreatedAt:     chapterModel.CreatedAt.Format(utils.TimeLayout),
		UpdatedAt:     chapterModel.UpdatedAt.Format(utils.TimeLayout),
	}, nil
}

func (s *GRPCChapterService) CreateChapter(ctx context.Context, req *chapter.CreateChapterRequest) (*chapter.CreateChapterResponse, error) {
	chapterRepo := repository_impl.NewChapterRepositoryImpl(s.db)
	
	chapterModel := models.NewChapterModel(
		req.MangaId,
		req.ChapterNumber,
		req.Title,
		req.PagesData,
	)
	if req.Id != "" {
		chapterModel.ID = req.Id
	}

	err := chapterRepo.SaveChapter(chapterModel)
	if err != nil {
		return nil, err
	}

	return &chapter.CreateChapterResponse{
		Success: true,
		Id:      chapterModel.ID,
	}, nil
}

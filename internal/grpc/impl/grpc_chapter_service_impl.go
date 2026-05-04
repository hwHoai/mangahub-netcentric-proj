package grpc_services_impl

import (
	"context"
	"mangahub/internal/grpc"
	"mangahub/pkg/utils"
	repository_impl "mangahub/pkg/repository/impl"
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

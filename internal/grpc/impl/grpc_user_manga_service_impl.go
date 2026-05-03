package grpc_services_impl

import (
	"context"
	"fmt"
	"mangahub/internal/grpc"
	"mangahub/pkg/models"
	"mangahub/pkg/models/enums"
	repository_impl "mangahub/pkg/repository/impl"
	manga "mangahub/proto/manga"
	"mangahub/proto/user_manga"
	"strings"

	"gorm.io/gorm"
)

type GRPCUserMangaService struct {
	user_manga.UnimplementedGRPCUserMangaServiceServer
	db *gorm.DB
}

var _ grpc.GRPCUserMangaService = (*GRPCUserMangaService)(nil)

func NewGRPCUserMangaService(db *gorm.DB) *GRPCUserMangaService {
	return &GRPCUserMangaService{db: db}
}

func (s *GRPCUserMangaService) FollowManga(ctx context.Context, req *user_manga.FollowMangaRequest) (*user_manga.FollowMangaResponse, error) {
	if req.UserId == "" || req.MangaId == "" {
		return nil, fmt.Errorf("user_id and manga_id are required")
	}

	followerRepo := repository_impl.NewMangaFollowerRepository(s.db)

	follower, err := followerRepo.FollowManga(req.UserId, req.MangaId)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("already following this manga")
		}
		return nil, fmt.Errorf("failed to follow manga: %w", err)
	}

	return &user_manga.FollowMangaResponse{
		UserId:    follower.UserModelID,
		MangaId:   follower.MangaModelID,
		CreatedAt: follower.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *GRPCUserMangaService) UnfollowManga(ctx context.Context, req *user_manga.UnfollowMangaRequest) (*user_manga.UnfollowMangaResponse, error) {
	if req.UserId == "" || req.MangaId == "" {
		return nil, fmt.Errorf("user_id and manga_id are required")
	}

	followerRepo := repository_impl.NewMangaFollowerRepository(s.db)
	err := followerRepo.UnfollowManga(req.UserId, req.MangaId)
	if err != nil {
		return nil, fmt.Errorf("failed to unfollow manga: %w", err)
	}

	return &user_manga.UnfollowMangaResponse{
		Success: true,
	}, nil
}

func (s *GRPCUserMangaService) GetFollowingMangas(ctx context.Context, req *user_manga.GetFollowingMangasRequest) (*user_manga.GetFollowingMangasResponse, error) {
	if req.UserId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	followerRepo := repository_impl.NewMangaFollowerRepository(s.db)
	mangas, err := followerRepo.GetFollowingMangas(req.UserId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, fmt.Errorf("failed to get following mangas: %w", err)
	}

	var mangaList []*manga.Manga
	for _, m := range mangas {
		mangaList = append(mangaList, &manga.Manga{
			Id:            m.ID,
			Title:         m.Title,
			Author:        m.Author,
			TotalChapters: int32(m.TotalChapters),
			Description:   m.Description,
			CoverUrl:      m.CoverURL,
			Status:        string(m.Status),
		})
	}

	return &user_manga.GetFollowingMangasResponse{
		Mangas: mangaList,
	}, nil
}

func (s *GRPCUserMangaService) StoreReadingProgress(ctx context.Context, req *user_manga.StoreReadingProgressRequest) (*user_manga.StoreReadingProgressResponse, error) {
	if req.UserId == "" || req.ChapterId == "" {
		return nil, fmt.Errorf("user_id and chapter_id are required")
	}

	// 1. Look up the chapter to get its manga_id and chapter_number
	chapterRepo := repository_impl.NewChapterRepositoryImpl(s.db)
	chapter, err := chapterRepo.GetChapterByID(req.ChapterId)
	if err != nil {
		return nil, fmt.Errorf("chapter not found: %w", err)
	}

	// 2. Upsert reading progress
	readingProgressRepo := repository_impl.NewReadingProgressRepository(s.db)
	progress := models.NewReadingProgressModel(
		req.UserId,
		chapter.MangaID,
		enums.ReadingStatusInProgress,
		int(chapter.ChapterNumber),
	)

	savedProgress, err := readingProgressRepo.UpsertReadingProgress(progress)
	if err != nil {
		return nil, fmt.Errorf("failed to store reading progress: %w", err)
	}

	return &user_manga.StoreReadingProgressResponse{
		UserId:         savedProgress.UserID,
		MangaId:        savedProgress.MangaID,
		Status:         string(savedProgress.Status),
		CurrentChapter: int32(savedProgress.CurrentChapter),
		UpdatedAt:      savedProgress.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

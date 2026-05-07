package grpc_services_impl

import (
	"context"
	"mangahub/proto/user_manga"

	"fmt"
	"mangahub/internal/grpc"
	"mangahub/internal/repository"
	repository_impl "mangahub/internal/repository/impl"
	"mangahub/pkg/models"
	"mangahub/pkg/models/enums"
	"mangahub/pkg/utils"
	manga "mangahub/proto/manga"
	"strings"

	"gorm.io/gorm"
)

type GRPCUserMangaService struct {
	user_manga.UnimplementedGRPCUserMangaServiceServer
	db              *gorm.DB
	followerRepo    repository.MangaFollowerRepository
	chapterRepo     repository.ChapterRepository
	progressRepo    repository.ReadingProgressRepository
}

var _ grpc.GRPCUserMangaService = (*GRPCUserMangaService)(nil)

func NewGRPCUserMangaService(db *gorm.DB) *GRPCUserMangaService {
	return &GRPCUserMangaService{
		db:           db,
		followerRepo: repository_impl.NewMangaFollowerRepository(db),
		chapterRepo:  repository_impl.NewChapterRepositoryImpl(db),
		progressRepo: repository_impl.NewReadingProgressRepository(db),
	}
}

func (s *GRPCUserMangaService) FollowManga(ctx context.Context, req *user_manga.FollowMangaRequest) (*user_manga.FollowMangaResponse, error) {
	if req.UserId == "" || req.MangaId == "" {
		return nil, fmt.Errorf("user_id and manga_id are required")
	}

	follower, err := s.followerRepo.FollowManga(req.UserId, req.MangaId)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("already following this manga")
		}
		return nil, fmt.Errorf("failed to follow manga: %w", err)
	}

	return &user_manga.FollowMangaResponse{
		UserId:    follower.UserModelID,
		MangaId:   follower.MangaModelID,
		CreatedAt: follower.CreatedAt.Format(utils.TimeLayout),
	}, nil
}

func (s *GRPCUserMangaService) UnfollowManga(ctx context.Context, req *user_manga.UnfollowMangaRequest) (*user_manga.UnfollowMangaResponse, error) {
	if req.UserId == "" || req.MangaId == "" {
		return nil, fmt.Errorf("user_id and manga_id are required")
	}

	err := s.followerRepo.UnfollowManga(req.UserId, req.MangaId)
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

	mangas, err := s.followerRepo.GetFollowingMangas(req.UserId, int(req.Limit), int(req.Offset))
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
			CreatedAt:     m.CreatedAt.Format(utils.TimeLayout),
			UpdatedAt:     m.UpdatedAt.Format(utils.TimeLayout),
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
	chapter, err := s.chapterRepo.GetChapterByID(req.ChapterId)
	if err != nil {
		return nil, fmt.Errorf("chapter not found: %w", err)
	}

	// 2. Upsert reading progress
	progress := models.NewReadingProgressModel(
		req.UserId,
		chapter.MangaID,
		enums.ReadingStatusInProgress,
		int(chapter.ChapterNumber),
	)

	savedProgress, err := s.progressRepo.UpsertReadingProgress(progress)
	if err != nil {
		return nil, fmt.Errorf("failed to store reading progress: %w", err)
	}

	return &user_manga.StoreReadingProgressResponse{
		UserId:         savedProgress.UserID,
		MangaId:        savedProgress.MangaID,
		Status:         string(savedProgress.Status),
		CurrentChapter: int32(savedProgress.CurrentChapter),
		UpdatedAt:      savedProgress.UpdatedAt.Format(utils.TimeLayout),
		CreatedAt:      savedProgress.CreatedAt.Format(utils.TimeLayout),
	}, nil
}
func (s *GRPCUserMangaService) GetReadingHistory(ctx context.Context, req *user_manga.GetReadingHistoryRequest) (*user_manga.GetReadingHistoryResponse, error) {
	if req.UserId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	history, err := s.progressRepo.GetReadingHistory(req.UserId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, fmt.Errorf("failed to get reading history: %w", err)
	}

	var historyItems []*user_manga.ReadingHistoryItem
	for _, item := range history {
		historyItems = append(historyItems, &user_manga.ReadingHistoryItem{
			Manga: &manga.Manga{
				Id:            item.Manga.ID,
				Title:         item.Manga.Title,
				Author:        item.Manga.Author,
				TotalChapters: int32(item.Manga.TotalChapters),
				Description:   item.Manga.Description,
				CoverUrl:      item.Manga.CoverURL,
				Status:        string(item.Manga.Status),
				CreatedAt:     item.Manga.CreatedAt.Format(utils.TimeLayout),
				UpdatedAt:     item.Manga.UpdatedAt.Format(utils.TimeLayout),
			},
			Status:         string(item.Status),
			CurrentChapter: int32(item.CurrentChapter),
			UpdatedAt:      item.UpdatedAt.Format(utils.TimeLayout),
			CreatedAt:      item.CreatedAt.Format(utils.TimeLayout),
		})
	}

	return &user_manga.GetReadingHistoryResponse{
		History: historyItems,
	}, nil
}

func (s *GRPCUserMangaService) GetFollowers(ctx context.Context, req *user_manga.GetFollowersRequest) (*user_manga.GetFollowersResponse, error) {
	if req.MangaId == "" {
		return nil, fmt.Errorf("manga_id is required")
	}

	userIDs, err := s.followerRepo.GetFollowers(req.MangaId)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	return &user_manga.GetFollowersResponse{
		UserIds: userIDs,
	}, nil
}

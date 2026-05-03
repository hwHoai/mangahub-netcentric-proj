package user_services_impl

import (
	"context"
	"fmt"
	user_services "mangahub/internal/user"
	"mangahub/pkg/models"
	"mangahub/pkg/models/enums"
	"mangahub/proto/user_manga"
)

type UserMangaServiceImpl struct {
	grpcUserMangaClient user_manga.GRPCUserMangaServiceClient
}

var _ user_services.UserMangaService = (*UserMangaServiceImpl)(nil)

func NewUserMangaService(grpcUserMangaClient user_manga.GRPCUserMangaServiceClient) user_services.UserMangaService {
	return &UserMangaServiceImpl{grpcUserMangaClient: grpcUserMangaClient}
}

func (s *UserMangaServiceImpl) FollowManga(userID string, mangaID string) (*models.MangaFollowerModel, error) {
	grpcRequest := &user_manga.FollowMangaRequest{
		UserId:  userID,
		MangaId: mangaID,
	}
	grpcResponse, err := s.grpcUserMangaClient.FollowManga(context.Background(), grpcRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to follow manga: %w", err)
	}

	return &models.MangaFollowerModel{
		UserModelID:  grpcResponse.UserId,
		MangaModelID: grpcResponse.MangaId,
	}, nil
}

func (s *UserMangaServiceImpl) UnfollowManga(userID string, mangaID string) error {
	grpcRequest := &user_manga.UnfollowMangaRequest{
		UserId:  userID,
		MangaId: mangaID,
	}
	_, err := s.grpcUserMangaClient.UnfollowManga(context.Background(), grpcRequest)
	if err != nil {
		return fmt.Errorf("failed to unfollow manga: %w", err)
	}
	return nil
}

func (s *UserMangaServiceImpl) GetFollowingMangas(userID string, limit int, offset int) ([]models.MangaModel, error) {
	grpcRequest := &user_manga.GetFollowingMangasRequest{
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	grpcResponse, err := s.grpcUserMangaClient.GetFollowingMangas(context.Background(), grpcRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get following mangas: %w", err)
	}

	var mangas []models.MangaModel
	for _, m := range grpcResponse.Mangas {
		mangas = append(mangas, models.MangaModel{
			ID:            m.Id,
			Title:         m.Title,
			Author:        m.Author,
			TotalChapters: int(m.TotalChapters),
			Description:   m.Description,
			CoverURL:      m.CoverUrl,
			Status:        enums.MangaStatus(m.Status),
		})
	}
	return mangas, nil
}

func (s *UserMangaServiceImpl) StoreReadingProgress(userID string, chapterID string) (*models.ReadingProgressModel, error) {
	grpcRequest := &user_manga.StoreReadingProgressRequest{
		UserId:    userID,
		ChapterId: chapterID,
	}
	grpcResponse, err := s.grpcUserMangaClient.StoreReadingProgress(context.Background(), grpcRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to store reading progress: %w", err)
	}

	return &models.ReadingProgressModel{
		UserID:         grpcResponse.UserId,
		MangaID:        grpcResponse.MangaId,
		Status:         enums.ReadingStatus(grpcResponse.Status),
		CurrentChapter: int(grpcResponse.CurrentChapter),
	}, nil
}

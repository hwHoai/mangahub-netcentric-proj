package user_services_impl

import (
	"context"
	"fmt"
	user_services "mangahub/internal/user"
	"mangahub/pkg/models"
	"mangahub/pkg/models/enums"
	"mangahub/pkg/utils"
	"mangahub/proto/user_manga"
	"time"
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

	createdAt, _ := time.Parse(utils.TimeLayout, grpcResponse.CreatedAt)

	return &models.MangaFollowerModel{
		UserModelID:  grpcResponse.UserId,
		MangaModelID: grpcResponse.MangaId,
		BaseModel: models.BaseModel{
			CreatedAt: createdAt,
		},
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
		createdAt, _ := time.Parse(utils.TimeLayout, m.CreatedAt)
		updatedAt, _ := time.Parse(utils.TimeLayout, m.UpdatedAt)

		mangas = append(mangas, models.MangaModel{
			ID:            m.Id,
			Title:         m.Title,
			Author:        m.Author,
			TotalChapters: int(m.TotalChapters),
			Description:   m.Description,
			CoverURL:      m.CoverUrl,
			Status:        enums.MangaStatus(m.Status),
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

func (s *UserMangaServiceImpl) StoreReadingProgress(userID string, chapterID string) (*models.ReadingProgressModel, error) {
	grpcRequest := &user_manga.StoreReadingProgressRequest{
		UserId:    userID,
		ChapterId: chapterID,
	}
	grpcResponse, err := s.grpcUserMangaClient.StoreReadingProgress(context.Background(), grpcRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to store reading progress: %w", err)
	}

	createdAt, _ := time.Parse(utils.TimeLayout, grpcResponse.CreatedAt)
	updatedAt, _ := time.Parse(utils.TimeLayout, grpcResponse.UpdatedAt)
	return &models.ReadingProgressModel{
		UserID:         grpcResponse.UserId,
		MangaID:        grpcResponse.MangaId,
		Status:         enums.ReadingStatus(grpcResponse.Status),
		CurrentChapter: int(grpcResponse.CurrentChapter),
		BaseModel: models.BaseModel{
			CreatedAt: createdAt,
		},
		MetaUpdateModel: models.MetaUpdateModel{
			UpdatedAt: updatedAt,
		},
	}, nil
}
func (s *UserMangaServiceImpl) GetReadingHistory(userID string, limit int, offset int) ([]models.ReadingProgressModel, error) {
	grpcRequest := &user_manga.GetReadingHistoryRequest{
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	grpcResponse, err := s.grpcUserMangaClient.GetReadingHistory(context.Background(), grpcRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get reading history: %w", err)
	}

	var history []models.ReadingProgressModel
	for _, item := range grpcResponse.History {
		updatedAt, _ := time.Parse(utils.TimeLayout, item.UpdatedAt)
		createdAt, _ := time.Parse(utils.TimeLayout, item.CreatedAt)

		mangaCreatedAt, _ := time.Parse(utils.TimeLayout, item.Manga.CreatedAt)
		mangaUpdatedAt, _ := time.Parse(utils.TimeLayout, item.Manga.UpdatedAt)

		history = append(history, models.ReadingProgressModel{
			UserID:  userID,
			MangaID: item.Manga.Id,
			Manga: models.MangaModel{
				ID:            item.Manga.Id,
				Title:         item.Manga.Title,
				Author:        item.Manga.Author,
				TotalChapters: int(item.Manga.TotalChapters),
				Description:   item.Manga.Description,
				CoverURL:      item.Manga.CoverUrl,
				Status:        enums.MangaStatus(item.Manga.Status),
				BaseModel: models.BaseModel{
					CreatedAt: mangaCreatedAt,
				},
				MetaUpdateModel: models.MetaUpdateModel{
					UpdatedAt: mangaUpdatedAt,
				},
			},
			Status:         enums.ReadingStatus(item.Status),
			CurrentChapter: int(item.CurrentChapter),
			BaseModel: models.BaseModel{
				CreatedAt: createdAt,
			},
			MetaUpdateModel: models.MetaUpdateModel{
				UpdatedAt: updatedAt,
			},
		})
	}
	return history, nil
}

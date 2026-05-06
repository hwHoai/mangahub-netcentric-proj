package manga_services_impl

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mangahub/internal/manga"
	"mangahub/pkg/models"
	"mangahub/pkg/utils"
	"mangahub/proto/chapter"
	"strconv"
	"time"

	udp_services "mangahub/internal/udp"
	"mangahub/pkg/clients"
)

type ChapterServiceImpl struct {
	grpcChapterClient       chapter.GRPCChapterServiceClient
	udpNotificationServices udp_services.UDPChapterNotificationServices
	mangaDexClient          clients.MangaDexClientInterface
}

func NewChapterService(
	grpcChapterClient chapter.GRPCChapterServiceClient,
	udpNotificationServices udp_services.UDPChapterNotificationServices,
	mangaDexClient clients.MangaDexClientInterface,
) manga.ChapterService {
	return &ChapterServiceImpl{
		grpcChapterClient:       grpcChapterClient,
		udpNotificationServices: udpNotificationServices,
		mangaDexClient:          mangaDexClient,
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

func (s *ChapterServiceImpl) CreateNewChapter(ctx context.Context, mangaID, mangadexChapterID string) (string, error) {
	// 1. Fetch data from MangaDex
	chapterDetails, err := s.mangaDexClient.GetChapterDetails(mangadexChapterID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch chapter details from MangaDex: %w", err)
	}

	pages, err := s.mangaDexClient.GetChapterPages(mangadexChapterID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch chapter pages from MangaDex: %w", err)
	}

	// 2. Parse details
	chapterNumber := 0.0
	if chapterDetails.Attributes.Chapter != "" {
		parsed, err := strconv.ParseFloat(chapterDetails.Attributes.Chapter, 64)
		if err == nil {
			chapterNumber = parsed
		}
	}
	pagesDataBytes, _ := json.Marshal(pages)

	// 3. Save via gRPC
	createReq := &chapter.CreateChapterRequest{
		Id:            chapterDetails.ID,
		MangaId:       mangaID,
		ChapterNumber: chapterNumber,
		Title:         chapterDetails.Attributes.Title,
		PagesData:     string(pagesDataBytes),
	}

	createRes, err := s.grpcChapterClient.CreateChapter(ctx, createReq)
	if err != nil {
		return "", fmt.Errorf("failed to save chapter: %w", err)
	}

	// 4. Send UDP Notification
	if s.udpNotificationServices != nil {
		s.udpNotificationServices.SendNewChapterNotification(mangaID, createRes.Id, chapterDetails.Attributes.Title, chapterNumber)
	} else {
		log.Printf("UDP Notification Services is not initialized")
	}

	return createRes.Id, nil
}

package manga_services_impl

import (
	"context"
	"encoding/json"
	"testing"

	"mangahub/pkg/clients"
	"mangahub/proto/chapter"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// --- Mocks ---

type MockGRPCChapterServiceClient struct {
	mock.Mock
}

func (m *MockGRPCChapterServiceClient) GetChapterByID(ctx context.Context, in *chapter.GetChapterByIDRequest, opts ...grpc.CallOption) (*chapter.GetChapterByIDResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*chapter.GetChapterByIDResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGRPCChapterServiceClient) CreateChapter(ctx context.Context, in *chapter.CreateChapterRequest, opts ...grpc.CallOption) (*chapter.CreateChapterResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*chapter.CreateChapterResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockUDPChapterNotificationServices struct {
	mock.Mock
}

func (m *MockUDPChapterNotificationServices) SendNewChapterNotification(mangaID string, chapterID string, title string, chapterNumber float64) error {
	args := m.Called(mangaID, chapterID, title, chapterNumber)
	return args.Error(0)
}

func (m *MockUDPChapterNotificationServices) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUDPChapterNotificationServices) SendNewMessageNotification(roomID, senderName, content string) error {
	args := m.Called(roomID, senderName, content)
	return args.Error(0)
}

type MockMangaDexClient struct {
	mock.Mock
}

func (m *MockMangaDexClient) GetChapterDetails(mangadexChapterID string) (*clients.MangaDexChapter, error) {
	args := m.Called(mangadexChapterID)
	if args.Get(0) != nil {
		return args.Get(0).(*clients.MangaDexChapter), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMangaDexClient) GetChapterPages(mangadexChapterID string) ([]string, error) {
	args := m.Called(mangadexChapterID)
	if args.Get(0) != nil {
		return args.Get(0).([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMangaDexClient) GetMangaList(limit int, offset int) (*clients.MangaDexListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).(*clients.MangaDexListResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMangaDexClient) GetTags() (map[string]clients.TagAttributes, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(map[string]clients.TagAttributes), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMangaDexClient) GetMangaChapters(mangaID string, limit int) (*clients.MangaDexChapterListResponse, error) {
	args := m.Called(mangaID, limit)
	if args.Get(0) != nil {
		return args.Get(0).(*clients.MangaDexChapterListResponse), args.Error(1)
	}
	return nil, args.Error(1)
}


// --- Tests ---

func TestReadChapter_Success(t *testing.T) {
	mockClient := new(MockGRPCChapterServiceClient)

	mockClient.On("GetChapterByID", mock.Anything, &chapter.GetChapterByIDRequest{
		ChapterId: "chapter-1",
	}).Return(&chapter.GetChapterByIDResponse{
		Id:            "chapter-1",
		MangaId:       "manga-1",
		ChapterNumber: 1,
		Title:         "Romance Dawn",
		PagesData:     "[\"page1.jpg\"]",
		CreatedAt:     "2026-01-01T00:00:00Z",
		UpdatedAt:     "2026-01-01T00:00:00Z",
	}, nil)

	service := NewChapterService(mockClient, nil, nil)

	chapterModel, err := service.ReadChapter("chapter-1")

	assert.NoError(t, err)
	assert.NotNil(t, chapterModel)
	assert.Equal(t, "chapter-1", chapterModel.ID)
	assert.Equal(t, "manga-1", chapterModel.MangaID)
	assert.Equal(t, 1.0, chapterModel.ChapterNumber)
	assert.Equal(t, "Romance Dawn", chapterModel.Title)

	mockClient.AssertExpectations(t)
}

func TestCreateNewChapter_Success(t *testing.T) {
	mockChapterClient := new(MockGRPCChapterServiceClient)
	mockUDPClient := new(MockUDPChapterNotificationServices)
	mockMangaDexClient := new(MockMangaDexClient)

	// Mock MangaDex response
	mockMangaDexClient.On("GetChapterDetails", "md-chapter-1").Return(&clients.MangaDexChapter{
		ID: "md-chapter-1",
		Attributes: clients.MangaDexChapterAttribute{
			Title:   "Romance Dawn",
			Chapter: "1",
			Pages:   2,
		},
	}, nil)

	pages := []string{"url1", "url2"}
	mockMangaDexClient.On("GetChapterPages", "md-chapter-1").Return(pages, nil)
	pagesDataBytes, _ := json.Marshal(pages)

	// Mock gRPC response
	mockChapterClient.On("CreateChapter", mock.Anything, &chapter.CreateChapterRequest{
		Id:            "md-chapter-1",
		MangaId:       "manga-1",
		ChapterNumber: 1.0,
		Title:         "Romance Dawn",
		PagesData:     string(pagesDataBytes),
	}).Return(&chapter.CreateChapterResponse{
		Id: "md-chapter-1",
	}, nil)

	// Mock UDP Notification
	mockUDPClient.On("SendNewChapterNotification", "manga-1", "md-chapter-1", "Romance Dawn", 1.0).Return(nil)
	mockUDPClient.On("Close").Return(nil).Maybe()

	service := NewChapterService(mockChapterClient, mockUDPClient, mockMangaDexClient)

	id, err := service.CreateNewChapter(context.Background(), "manga-1", "md-chapter-1")

	assert.NoError(t, err)
	assert.Equal(t, "md-chapter-1", id)

	mockChapterClient.AssertExpectations(t)
	mockUDPClient.AssertExpectations(t)
	mockMangaDexClient.AssertExpectations(t)
}

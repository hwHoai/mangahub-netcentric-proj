package manga_services_impl

import (
	"context"
	"testing"

	"mangahub/proto/manga"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// --- Mocks ---

type MockGRPCMangaServiceClient struct {
	mock.Mock
}

func (m *MockGRPCMangaServiceClient) GetMangas(ctx context.Context, in *manga.MangaListRequest, opts ...grpc.CallOption) (*manga.MangaListResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*manga.MangaListResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGRPCMangaServiceClient) GetMangaDetail(ctx context.Context, in *manga.MangaDetailRequest, opts ...grpc.CallOption) (*manga.MangaDetailResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*manga.MangaDetailResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGRPCMangaServiceClient) GetMangaChapters(ctx context.Context, in *manga.MangaChaptersRequest, opts ...grpc.CallOption) (*manga.MangaChaptersResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*manga.MangaChaptersResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGRPCMangaServiceClient) CheckMangaExists(ctx context.Context, in *manga.CheckMangaExistsRequest, opts ...grpc.CallOption) (*manga.CheckMangaExistsResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*manga.CheckMangaExistsResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// --- Tests ---

func TestGetMangaDetail_Success(t *testing.T) {
	mockClient := new(MockGRPCMangaServiceClient)

	mockClient.On("GetMangaDetail", mock.Anything, &manga.MangaDetailRequest{
		Id: "manga-1",
	}).Return(&manga.MangaDetailResponse{
		Manga: &manga.Manga{
			Id:            "manga-1",
			Title:         "One Piece",
			Author:        "Eiichiro Oda",
			TotalChapters: 1050,
			Status:        "1",
			CreatedAt:     "2026-01-01T00:00:00Z",
			UpdatedAt:     "2026-01-01T00:00:00Z",
		},
	}, nil)

	service := NewMangaService(mockClient)

	mangaModel, err := service.GetMangaDetail("manga-1")

	assert.NoError(t, err)
	assert.NotNil(t, mangaModel)
	assert.Equal(t, "manga-1", mangaModel.ID)
	assert.Equal(t, "One Piece", mangaModel.Title)
	assert.Equal(t, "Eiichiro Oda", mangaModel.Author)

	mockClient.AssertExpectations(t)
}

func TestListMangas_Success(t *testing.T) {
	mockClient := new(MockGRPCMangaServiceClient)

	mockClient.On("GetMangas", mock.Anything, &manga.MangaListRequest{
		Limit:  10,
		Offset: 0,
	}).Return(&manga.MangaListResponse{
		Mangas: []*manga.Manga{
			{
				Id:        "manga-1",
				Title:     "One Piece",
				CreatedAt: "2026-01-01T00:00:00Z",
				UpdatedAt: "2026-01-01T00:00:00Z",
			},
			{
				Id:        "manga-2",
				Title:     "Naruto",
				CreatedAt: "2026-01-01T00:00:00Z",
				UpdatedAt: "2026-01-01T00:00:00Z",
			},
		},
	}, nil)

	service := NewMangaService(mockClient)

	mangas, err := service.ListMangas(10, 0)

	assert.NoError(t, err)
	assert.Len(t, mangas, 2)
	assert.Equal(t, "One Piece", mangas[0].Title)
	assert.Equal(t, "Naruto", mangas[1].Title)

	mockClient.AssertExpectations(t)
}

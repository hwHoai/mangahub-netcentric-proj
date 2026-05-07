package user_services_impl

import (
	"context"
	"testing"

	"mangahub/proto/user_manga"
	"mangahub/proto/manga"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// --- Mocks ---

type MockGRPCUserMangaServiceClient struct {
	mock.Mock
}

func (m *MockGRPCUserMangaServiceClient) FollowManga(ctx context.Context, in *user_manga.FollowMangaRequest, opts ...grpc.CallOption) (*user_manga.FollowMangaResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*user_manga.FollowMangaResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGRPCUserMangaServiceClient) UnfollowManga(ctx context.Context, in *user_manga.UnfollowMangaRequest, opts ...grpc.CallOption) (*user_manga.UnfollowMangaResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*user_manga.UnfollowMangaResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGRPCUserMangaServiceClient) GetFollowingMangas(ctx context.Context, in *user_manga.GetFollowingMangasRequest, opts ...grpc.CallOption) (*user_manga.GetFollowingMangasResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*user_manga.GetFollowingMangasResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGRPCUserMangaServiceClient) StoreReadingProgress(ctx context.Context, in *user_manga.StoreReadingProgressRequest, opts ...grpc.CallOption) (*user_manga.StoreReadingProgressResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*user_manga.StoreReadingProgressResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGRPCUserMangaServiceClient) GetReadingHistory(ctx context.Context, in *user_manga.GetReadingHistoryRequest, opts ...grpc.CallOption) (*user_manga.GetReadingHistoryResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*user_manga.GetReadingHistoryResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGRPCUserMangaServiceClient) GetFollowers(ctx context.Context, in *user_manga.GetFollowersRequest, opts ...grpc.CallOption) (*user_manga.GetFollowersResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*user_manga.GetFollowersResponse), args.Error(1)
	}
	return nil, args.Error(1)
}


// --- Tests ---

func TestFollowManga_Success(t *testing.T) {
	mockClient := new(MockGRPCUserMangaServiceClient)

	mockClient.On("FollowManga", mock.Anything, &user_manga.FollowMangaRequest{
		UserId:  "user-1",
		MangaId: "manga-1",
	}).Return(&user_manga.FollowMangaResponse{
		UserId:    "user-1",
		MangaId:   "manga-1",
		CreatedAt: "2026-01-01T00:00:00Z",
	}, nil)

	service := NewUserMangaService(mockClient)

	resp, err := service.FollowManga("user-1", "manga-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "user-1", resp.UserModelID)
	assert.Equal(t, "manga-1", resp.MangaModelID)

	mockClient.AssertExpectations(t)
}

func TestGetFollowingMangas_Success(t *testing.T) {
	mockClient := new(MockGRPCUserMangaServiceClient)

	mockClient.On("GetFollowingMangas", mock.Anything, &user_manga.GetFollowingMangasRequest{
		UserId: "user-1",
		Limit:  10,
		Offset: 0,
	}).Return(&user_manga.GetFollowingMangasResponse{
		Mangas: []*manga.Manga{
			{
				Id:            "manga-1",
				Title:         "Naruto",
				TotalChapters: 700,
				Status:        "1",
				CreatedAt:     "2026-01-01T00:00:00Z",
				UpdatedAt:     "2026-01-01T00:00:00Z",
			},
		},
	}, nil)

	service := NewUserMangaService(mockClient)

	mangas, err := service.GetFollowingMangas("user-1", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, mangas, 1)
	assert.Equal(t, "manga-1", mangas[0].ID)
	assert.Equal(t, "Naruto", mangas[0].Title)

	mockClient.AssertExpectations(t)
}

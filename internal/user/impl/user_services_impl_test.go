package user_services_impl

import (
	"context"
	"errors"
	"testing"
	
	user_proto "mangahub/proto/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// --- Mocks ---

type MockGRPCUserServiceClient struct {
	mock.Mock
}

func (m *MockGRPCUserServiceClient) GetUserModelByUsername(ctx context.Context, in *user_proto.GetUserModelByUsernameRequest, opts ...grpc.CallOption) (*user_proto.GetUserModelByUsernameResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*user_proto.GetUserModelByUsernameResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGRPCUserServiceClient) CreateNewUser(ctx context.Context, in *user_proto.CreateNewUserRequest, opts ...grpc.CallOption) (*user_proto.CreateNewUserResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*user_proto.CreateNewUserResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGRPCUserServiceClient) GetUserByID(ctx context.Context, in *user_proto.GetUserByIDRequest, opts ...grpc.CallOption) (*user_proto.GetUserByIDResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*user_proto.GetUserByIDResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// --- Tests ---

func TestGetUserDetails_Success(t *testing.T) {
	mockUserClient := new(MockGRPCUserServiceClient)

	// Setup mock behavior
	mockUserClient.On("GetUserByID", mock.Anything, &user_proto.GetUserByIDRequest{
		UserId: "user-123",
	}).Return(&user_proto.GetUserByIDResponse{
		UserId:    "user-123",
		Username:  "testuser",
		CreatedAt: "2026-01-01T00:00:00Z",
		UpdatedAt: "2026-01-01T00:00:00Z",
	}, nil)

	service := NewUserService(mockUserClient)

	resp, exception := service.GetUserDetails("user-123")

	assert.Equal(t, 0, exception.Code)
	assert.NotNil(t, resp)
	assert.Equal(t, "user-123", resp.UserID)
	assert.Equal(t, "testuser", resp.Username)

	mockUserClient.AssertExpectations(t)
}

func TestGetUserDetails_EmptyID(t *testing.T) {
	service := NewUserService(nil) // Client not needed for this path

	resp, exception := service.GetUserDetails("")

	assert.Nil(t, resp)
	assert.Equal(t, 401, exception.Code)
	assert.Equal(t, "Unauthorized", exception.Message)
}

func TestGetUserDetails_NotFound(t *testing.T) {
	mockUserClient := new(MockGRPCUserServiceClient)

	mockUserClient.On("GetUserByID", mock.Anything, &user_proto.GetUserByIDRequest{
		UserId: "non-existent-user",
	}).Return(nil, errors.New("user not found"))

	service := NewUserService(mockUserClient)

	resp, exception := service.GetUserDetails("non-existent-user")

	assert.Nil(t, resp)
	assert.Equal(t, 404, exception.Code)
	assert.Equal(t, "User not found", exception.Message)

	mockUserClient.AssertExpectations(t)
}

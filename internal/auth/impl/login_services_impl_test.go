package auth_service_impl

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"mangahub/pkg/dto"
	"mangahub/proto/session"
	"mangahub/proto/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// --- Mocks ---

type MockGRPCUserServiceClient struct {
	mock.Mock
}

func (m *MockGRPCUserServiceClient) GetUserModelByUsername(ctx context.Context, in *user.GetUserModelByUsernameRequest, opts ...grpc.CallOption) (*user.GetUserModelByUsernameResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*user.GetUserModelByUsernameResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGRPCUserServiceClient) CreateNewUser(ctx context.Context, in *user.CreateNewUserRequest, opts ...grpc.CallOption) (*user.CreateNewUserResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*user.CreateNewUserResponse), args.Error(1)
}

func (m *MockGRPCUserServiceClient) GetUserByID(ctx context.Context, in *user.GetUserByIDRequest, opts ...grpc.CallOption) (*user.GetUserByIDResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*user.GetUserByIDResponse), args.Error(1)
}

type MockGRPCSessionServiceClient struct {
	mock.Mock
}

func (m *MockGRPCSessionServiceClient) SaveSession(ctx context.Context, in *session.SaveSessionRequest, opts ...grpc.CallOption) (*session.SaveSessionResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*session.SaveSessionResponse), args.Error(1)
}

func (m *MockGRPCSessionServiceClient) UpdateSession(ctx context.Context, in *session.UpdateSessionRequest, opts ...grpc.CallOption) (*session.UpdateSessionResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*session.UpdateSessionResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGRPCSessionServiceClient) GetSessionByRefreshToken(ctx context.Context, in *session.GetSessionByRefreshTokenRequest, opts ...grpc.CallOption) (*session.GetSessionByRefreshTokenResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*session.GetSessionByRefreshTokenResponse), args.Error(1)
}

// --- Tests ---

func TestLoginByUsername_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup Mocks
	mockUserClient := new(MockGRPCUserServiceClient)
	mockSessionClient := new(MockGRPCSessionServiceClient)

	// Setup crypto
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	password := "mysecurepassword"
	hash := sha256.Sum256([]byte(password))
	hashedPassword := fmt.Sprintf("%x", hash)

	// Expected mock calls
	mockUserClient.On("GetUserModelByUsername", mock.Anything, &user.GetUserModelByUsernameRequest{Username: "testuser"}).
		Return(&user.GetUserModelByUsernameResponse{
			UserId:         "user-123",
			Username:       "testuser",
			HashedPassword: hashedPassword,
		}, nil)

	mockSessionClient.On("UpdateSession", mock.Anything, mock.AnythingOfType("*session.UpdateSessionRequest")).
		Return(&session.UpdateSessionResponse{}, nil)

	// Setup Gin Context
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest("POST", "/login", nil)

	// Initialize service
	service := &LoginServiceImpl{
		GRPCUserClient:    mockUserClient,
		GRPCSessionClient: mockSessionClient,
		Context:           ctx,
		PrivateKey:        privateKey,
	}

	// Execute
	req := &dto.LoginByUsernameRequest{
		Username: "testuser",
		Password: password,
	}
	resp, exception := service.LoginByUsername(req)

	// Assert
	assert.Equal(t, 0, exception.Code)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	mockUserClient.AssertExpectations(t)
	mockSessionClient.AssertExpectations(t)
}

func TestLoginByUsername_InvalidPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserClient := new(MockGRPCUserServiceClient)

	// Return a user but with a different hashed password
	mockUserClient.On("GetUserModelByUsername", mock.Anything, mock.Anything).
		Return(&user.GetUserModelByUsernameResponse{
			UserId:         "user-123",
			Username:       "testuser",
			HashedPassword: "different_hash",
		}, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest("POST", "/login", nil)

	service := &LoginServiceImpl{
		GRPCUserClient: mockUserClient,
		Context:        ctx,
		PrivateKey:     nil, // Not needed for this test as it fails before signing
	}

	req := &dto.LoginByUsernameRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	resp, exception := service.LoginByUsername(req)

	assert.Nil(t, resp)
	assert.Equal(t, 401, exception.Code)
	assert.Equal(t, "Invalid username or password", exception.Message)
}

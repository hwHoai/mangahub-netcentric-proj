package grpc_services_impl

import (
	"context"
	"errors"
	"testing"
	"time"

	"mangahub/pkg/models"
	"mangahub/proto/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByUsername(username string) (*models.UserModel, error) {
	args := m.Called(username)
	if args.Get(0) != nil {
		return args.Get(0).(*models.UserModel), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *models.UserModel) (*models.UserModel, error) {
	args := m.Called(user)
	if args.Get(0) != nil {
		return args.Get(0).(*models.UserModel), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id string) (*models.UserModel, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.UserModel), args.Error(1)
	}
	return nil, args.Error(1)
}

// --- Tests ---

func TestGRPCUserService_GetUserModelByUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)

	testTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	mockRepo.On("GetUserByUsername", "testuser").Return(&models.UserModel{
		ID:             "user-1",
		Username:       "testuser",
		HashedPassword: "hashed_password",
		BaseModel: models.BaseModel{
			CreatedAt: testTime,
		},
		MetaUpdateModel: models.MetaUpdateModel{
			UpdatedAt: testTime,
		},
	}, nil)

	service := &GRPCUserService{
		userRepo: mockRepo,
	}

	req := &user.GetUserModelByUsernameRequest{
		Username: "testuser",
	}

	resp, err := service.GetUserModelByUsername(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "user-1", resp.UserId)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "hashed_password", resp.HashedPassword)

	mockRepo.AssertExpectations(t)
}

func TestGRPCUserService_CreateNewUser(t *testing.T) {
	mockRepo := new(MockUserRepository)

	testTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	mockRepo.On("CreateUser", mock.AnythingOfType("*models.UserModel")).Return(&models.UserModel{
		ID:       "user-new",
		Username: "newuser",
		BaseModel: models.BaseModel{
			CreatedAt: testTime,
		},
		MetaUpdateModel: models.MetaUpdateModel{
			UpdatedAt: testTime,
		},
	}, nil)

	service := &GRPCUserService{
		userRepo: mockRepo,
	}

	req := &user.CreateNewUserRequest{
		Username: "newuser",
		Password: "password123",
	}

	resp, err := service.CreateNewUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "user-new", resp.UserId)
	assert.Equal(t, "newuser", resp.Username)

	mockRepo.AssertExpectations(t)
}

func TestGRPCUserService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetUserByID", "unknown-user").Return(nil, errors.New("record not found"))

	service := &GRPCUserService{
		userRepo: mockRepo,
	}

	req := &user.GetUserByIDRequest{
		UserId: "unknown-user",
	}

	resp, err := service.GetUserByID(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "record not found", err.Error())

	mockRepo.AssertExpectations(t)
}

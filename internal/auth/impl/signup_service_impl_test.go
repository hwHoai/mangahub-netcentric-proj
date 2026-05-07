package auth_service_impl

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"mangahub/pkg/dto"
	"mangahub/proto/session"
	"mangahub/proto/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignupByUsername_Success(t *testing.T) {
	// Setup Mocks
	mockUserClient := new(MockGRPCUserServiceClient)
	mockSessionClient := new(MockGRPCSessionServiceClient)

	// Setup crypto
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	// Expected mock calls
	mockUserClient.On("CreateNewUser", mock.Anything, &user.CreateNewUserRequest{
		Username: "newuser",
		Password: "newpassword",
	}).Return(&user.CreateNewUserResponse{
		UserId:   "user-999",
		Username: "newuser",
	}, nil)

	mockSessionClient.On("SaveSession", mock.Anything, mock.AnythingOfType("*session.SaveSessionRequest")).
		Return(&session.SaveSessionResponse{}, nil)

	// Initialize service
	service := &SignupServiceImpl{
		GRPCUserClient:    mockUserClient,
		GRPCSessionClient: mockSessionClient,
		Context:           context.Background(),
		PrivateKey:        privateKey,
	}

	// Execute
	req := &dto.SignupByUsernameRequest{
		Username: "newuser",
		Password: "newpassword",
	}
	resp, exception := service.SignupByUsername(req)

	// Assert
	assert.Equal(t, 0, exception.Code)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	mockUserClient.AssertExpectations(t)
	mockSessionClient.AssertExpectations(t)
}

func TestSignupByUsername_MissingFields(t *testing.T) {
	service := &SignupServiceImpl{}

	req := &dto.SignupByUsernameRequest{
		Username: "",
		Password: "",
	}
	resp, exception := service.SignupByUsername(req)

	assert.Nil(t, resp)
	assert.Equal(t, 400, exception.Code)
	assert.Equal(t, "Username and password are required", exception.Message)
}

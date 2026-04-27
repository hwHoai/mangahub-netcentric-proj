package auth_service_impl

import (
	"fmt"
	"mangahub/internal/auth"
	"mangahub/pkg/dto"
	"mangahub/proto/user"

	"github.com/gin-gonic/gin"
)

type LoginServiceImpl struct {
	GRPCUserClient user.GRPCUserServiceClient
	Context *gin.Context
}

var _ auth.LoginService = (*LoginServiceImpl)(nil)

func (s *LoginServiceImpl) LoginByUsername(request *dto.LoginByUsernameRequest) (*dto.LoginByUsernameResponse, dto.ExceptionDTO) {
	// 1. Check the request is valid
	if request.Username == "" || request.Password == "" {
		return nil, dto.ExceptionDTO{
			Code:    400,
			Message: "Username and password are required",
		}
	}

	// 2. Check the user exists and the password is correct
	grpcRequest := &user.GetUserModelByUsernameRequest{
		Username: request.Username,
	}

	grpcResponse, err := s.GRPCUserClient.GetUserModelByUsername(s.Context.Request.Context(), grpcRequest)
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    404,
			Message: "User not found",
		}
	}

	fmt.Printf("User found: %v\n", grpcResponse)

	return &dto.LoginByUsernameResponse{
		AccessToken: "dummy_access_token_for_" + grpcResponse.Username, // In real implementation, generate a JWT or similar token
		RefreshToken: "dummy_refresh_token_for_" + grpcResponse.Username, // In real implementation, generate a refresh token
		ExpiresIn: 3600, // Token expiration time in seconds (e.g., 1 hour)
	}, dto.ExceptionDTO{}
}
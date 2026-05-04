package auth_service_impl

import (
	"crypto/rsa"
	"crypto/sha256"
	"fmt"

	"mangahub/internal/auth"
	"mangahub/pkg/dto"
	"mangahub/pkg/utils/jwt"
	jwt_impl "mangahub/pkg/utils/jwt/impl"
	"mangahub/proto/session"
	"mangahub/proto/user"

	"github.com/gin-gonic/gin"
)

type LoginServiceImpl struct {
	GRPCUserClient    user.GRPCUserServiceClient
	GRPCSessionClient session.GRPCSessionServiceClient
	Context           *gin.Context
	PrivateKey        *rsa.PrivateKey
}

var _ auth.LoginService = (*LoginServiceImpl)(nil)

func (s *LoginServiceImpl) LoginByUsername(request *dto.LoginByUsernameRequest) (*dto.LoginByUsernameResponse, dto.ExceptionDTO) {
	if request.Username == "" || request.Password == "" {
		return nil, dto.ExceptionDTO{
			Code:    400,
			Message: "Username and password are required",
		}
	}

	grpcRequest := &user.GetUserModelByUsernameRequest{
		Username: request.Username,
	}

	grpcResponse, err := s.GRPCUserClient.GetUserModelByUsername(s.Context.Request.Context(), grpcRequest)
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    401,
			Message: "Invalid username or password",
		}
	}

	inputHash := sha256.Sum256([]byte(request.Password))
	hashedInputPassword := fmt.Sprintf("%x", inputHash)
	if grpcResponse.HashedPassword != hashedInputPassword {
		return nil, dto.ExceptionDTO{
			Code:    401,
			Message: "Invalid username or password",
		}
	}

	privateKey := s.PrivateKey
	if privateKey == nil {
		return nil, dto.ExceptionDTO{
			Code:    500,
			Message: "Private key is unavailable",
		}
	}

	jwtUtil := jwt_impl.NewJWTUtil(s.GRPCSessionClient)
	accessToken, err := jwtUtil.SignJWTToken(grpcResponse.UserId, jwt.AccessTokenTTL, privateKey)
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    500,
			Message: "Failed to generate access token",
		}
	}

	refreshToken, err := jwtUtil.SignJWTToken(grpcResponse.UserId, jwt.RefreshTokenTTL, privateKey)
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    500,
			Message: "Failed to generate refresh token",
		}
	}

	if s.GRPCSessionClient != nil {
		_, err = s.GRPCSessionClient.UpdateSession(s.Context.Request.Context(), &session.UpdateSessionRequest{
			UserId:       grpcResponse.UserId,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
		if err != nil {
			return nil, dto.ExceptionDTO{
				Code:    500,
				Message: "Failed to update session",
			}
		}
	}

	return &dto.LoginByUsernameResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(jwt.AccessTokenTTL.Seconds()),
	}, dto.ExceptionDTO{}
}
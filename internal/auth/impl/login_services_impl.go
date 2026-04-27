package auth_service_impl

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"mangahub/internal/auth"
	"mangahub/pkg/dto"
	"mangahub/proto/session"
	"mangahub/proto/user"

	"github.com/gin-gonic/gin"
)

type LoginServiceImpl struct {
	GRPCUserClient    user.GRPCUserServiceClient
	GRPCSessionClient session.GRPCSessionServiceClient
	Context           *gin.Context
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

	// 2. Retrieve user data from gRPC
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

	// 3. Initialize JWT service and create RSA key pair for token signing
	jwtService := NewJWTService()
	privateKey, _, err := jwtService.CreateRSAKeyPair(2048)
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    500,
			Message: "Failed to create RSA key pair",
		}
	}

	// 4. Sign tokens
	accessToken, err := jwtService.SignJWTToken(grpcResponse.UserId, auth.AccessTokenTTL, privateKey)
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    500,
			Message: "Failed to generate access token",
		}
	}

	refreshToken, err := jwtService.SignJWTToken(grpcResponse.UserId, auth.RefreshTokenTTL, privateKey)
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    500,
			Message: "Failed to generate refresh token",
		}
	}

	// 5. Export public key to PEM format
	publicKey := privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    500,
			Message: "Failed to export public key",
		}
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// 6. Update session via gRPC
	if s.GRPCSessionClient != nil {
		_, err = s.GRPCSessionClient.UpdateSession(s.Context.Request.Context(), &session.UpdateSessionRequest{
			UserId:       grpcResponse.UserId,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			PublicKey:    string(publicKeyPEM),
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
		ExpiresIn:    int64(auth.AccessTokenTTL.Seconds()),
	}, dto.ExceptionDTO{}
}
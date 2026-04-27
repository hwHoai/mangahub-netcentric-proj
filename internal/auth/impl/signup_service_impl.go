package auth_service_impl

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"mangahub/internal/auth"
	"mangahub/pkg/dto"
	"mangahub/proto/session"
	"mangahub/proto/user"

	"gorm.io/gorm"
)

type SignupServiceImpl struct {
	Context              any
	GRPCUserClient       user.GRPCUserServiceClient
	GRPCSessionClient    session.GRPCSessionServiceClient
}

var _ auth.SignupService = (*SignupServiceImpl)(nil)

func (s *SignupServiceImpl) SignupByUsername(request *dto.SignupByUsernameRequest) (*dto.SignupByUsernameResponse, dto.ExceptionDTO) {
	if request == nil {
		return nil, dto.ExceptionDTO{
			Code:    400,
			Message: "Invalid request",
		}
	}

	// 1. Validate request
	if request.Username == "" || request.Password == "" {
		return nil, dto.ExceptionDTO{
			Code:    400,
			Message: "Username and password are required",
		}
	}

	// 2. Create user via gRPC
	ctx, ok := s.Context.(context.Context)
	if !ok || ctx == nil {
		ctx = context.Background()
	}

	grpcRequest := &user.CreateNewUserRequest{
		Username: request.Username,
		Password: request.Password,
	}

	grpcResponse, err := s.GRPCUserClient.CreateNewUser(ctx, grpcRequest)
	if err != nil {
		// Check if user already exists (unique constraint violation)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, dto.ExceptionDTO{
				Code:    409,
				Message: "User already exists",
			}
		}
		// Check if the error message contains duplicate key info
		if err.Error() != "" && len(err.Error()) > 0 {
			return nil, dto.ExceptionDTO{
				Code:    409,
				Message: "User already exists",
			}
		}
		return nil, dto.ExceptionDTO{
			Code:    500,
			Message: "Failed to create user: " + err.Error(),
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

	// 6. Save session to database via gRPC
	saveSessionReq := &session.SaveSessionRequest{
		UserId:       grpcResponse.UserId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		PublicKey:    string(publicKeyPEM),
	}

	_, err = s.GRPCSessionClient.SaveSession(ctx, saveSessionReq)
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    500,
			Message: "Failed to save session: " + err.Error(),
		}
	}

	// 7. Calculate expiry time
	expiresIn := int64(auth.AccessTokenTTL.Seconds())

	return &dto.SignupByUsernameResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, dto.ExceptionDTO{}
}

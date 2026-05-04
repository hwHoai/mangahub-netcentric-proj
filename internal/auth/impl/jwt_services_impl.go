package auth_service_impl

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"mangahub/internal/auth"
	"mangahub/pkg/dto"
	"mangahub/proto/session"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTServiceImpl struct {
	GRPCSessionClient session.GRPCSessionServiceClient
}

var _ auth.JWTService = (*JWTServiceImpl)(nil)

func NewJWTService(grpcSessionClient session.GRPCSessionServiceClient) *JWTServiceImpl {
	return &JWTServiceImpl{
		GRPCSessionClient: grpcSessionClient,
	}
}

func (s *JWTServiceImpl) CreateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if bits < 2048 {
		bits = 2048
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

func (s *JWTServiceImpl) SignJWTToken(subject string, expiresIn time.Duration, privateKey *rsa.PrivateKey) (string, error) {
	if subject == "" {
		return "", errors.New("subject is required")
	}
	if privateKey == nil {
		return "", errors.New("private key is required")
	}
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": subject, // User ID
		"iat": now.Unix(),
		"exp": now.Add(expiresIn).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func (s *JWTServiceImpl) VerifyJWTToken(token string, publicKey *rsa.PublicKey) (*auth.JWTClaims, error) {
	if token == "" {
		return nil, errors.New("token is required")
	}
	if publicKey == nil {
		return nil, errors.New("public key is required")
	}

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != auth.JWTAlgorithm {
			return nil, fmt.Errorf("unexpected signing algorithm: %s", t.Method.Alg())
		}
		return publicKey, nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	subject, _ := claims["sub"].(string)
	if subject == "" {
		return nil, errors.New("invalid token claims")
	}

	expiresAt, err := numericClaimToInt64(claims["exp"])
	if err != nil {
		return nil, errors.New("invalid token claims")
	}

	issuedAt, err := numericClaimToInt64(claims["iat"])
	if err != nil {
		return nil, errors.New("invalid token claims")
	}

	if s.IsExpire(expiresAt) {
		return nil, errors.New("token expired")
	}

	return &auth.JWTClaims{
		Subject:   subject,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *JWTServiceImpl) IsExpire(expUnix int64) bool {
	return time.Now().Unix() >= expUnix
}

func (s *JWTServiceImpl) ParsePublicKeyPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
	if publicKeyPEM == "" {
		return nil, errors.New("public key is empty")
	}

	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to decode public key PEM")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not RSA format")
	}

	return publicKey, nil
}

func (s *JWTServiceImpl) RefreshToken(request *dto.RefreshTokenRequest, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) (*dto.RefreshTokenResponse, dto.ExceptionDTO) {
	if request.RefreshToken == "" {
		return nil, dto.ExceptionDTO{
			Code:    400,
			Message: "Refresh token is required",
		}
	}

	// 1. Verify the refresh token cryptographically
	claims, err := s.VerifyJWTToken(request.RefreshToken, publicKey)
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    401,
			Message: "Invalid refresh token",
		}
	}

	if s.GRPCSessionClient == nil {
		return nil, dto.ExceptionDTO{
			Code:    500,
			Message: "Session service is unavailable",
		}
	}

	// 2. Check if the refresh token exists in the database
	grpcResponse, err := s.GRPCSessionClient.GetSessionByRefreshToken(context.Background(), &session.GetSessionByRefreshTokenRequest{
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    401,
			Message: "Refresh token not found or expired",
		}
	}

	// 3. Double check user_id matches claims
	if grpcResponse.UserId != claims.Subject {
		return nil, dto.ExceptionDTO{
			Code:    401,
			Message: "Invalid refresh token claims",
		}
	}

	// 4. Generate new tokens
	newAccessToken, err := s.SignJWTToken(claims.Subject, auth.AccessTokenTTL, privateKey)
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    500,
			Message: "Failed to generate access token",
		}
	}

	newRefreshToken, err := s.SignJWTToken(claims.Subject, auth.RefreshTokenTTL, privateKey)
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    500,
			Message: "Failed to generate refresh token",
		}
	}

	// 5. Update session in database
	_, err = s.GRPCSessionClient.UpdateSession(context.Background(), &session.UpdateSessionRequest{
		UserId:       claims.Subject,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    500,
			Message: "Failed to update session in database",
		}
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, dto.ExceptionDTO{}
}

func numericClaimToInt64(value any) (int64, error) {
	switch v := value.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case json.Number:
		i, err := v.Int64()
		if err != nil {
			return 0, err
		}
		return i, nil
	default:
		return 0, errors.New("invalid numeric claim type")
	}
}

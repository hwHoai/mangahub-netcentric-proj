package auth

import (
	"crypto/rsa"
	"mangahub/pkg/dto"
	"time"
)

const (
	JWTAlgorithm     = "RS256"
	AccessTokenType  = "access_token"
	RefreshTokenType = "refresh_token"
)

var (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

type JWTClaims struct {
	Subject   string `json:"sub"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

type JWTService interface {
	CreateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error)
	SignJWTToken(subject string, expiresIn time.Duration, privateKey *rsa.PrivateKey) (string, error)
	VerifyJWTToken(token string, publicKey *rsa.PublicKey) (*JWTClaims, error)
	IsExpire(expUnix int64) bool
	ParsePublicKeyPEM(publicKeyPEM string) (*rsa.PublicKey, error)
	RefreshToken(request *dto.RefreshTokenRequest, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) (*dto.RefreshTokenResponse, dto.ExceptionDTO)
}

package auth_service_impl

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"mangahub/internal/auth"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTServiceImpl struct{}

var _ auth.JWTService = (*JWTServiceImpl)(nil)

func NewJWTService() *JWTServiceImpl {
	return &JWTServiceImpl{}
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
		"sub":        subject, // User ID
		"iat":        now.Unix(),
		"exp":        now.Add(expiresIn).Unix(),
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

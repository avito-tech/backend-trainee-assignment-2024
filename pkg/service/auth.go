package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthConfig struct {
	SigningKey string
	TokenTTL   time.Duration
}

type AuthService struct {
	config AuthConfig
}

func NewAuthService(config AuthConfig) *AuthService {
	return &AuthService{config: config}
}

type tokenClaims struct {
	jwt.StandardClaims
	Role string `json:"role"`
}

func (s *AuthService) GenerateToken(role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.config.TokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		role,
	})

	return token.SignedString([]byte(s.config.SigningKey))
}

package service

import (
	"errors"
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

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.config.SigningKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.Role, nil
}

package utils

import (
	"net/http"
	"os"
	"testing"

	"github.com/spf13/viper"

	"gta2024/pkg/service"
)

func GenerateToken(t *testing.T, role string) string {
	authConfig := service.AuthConfig{
		TokenTTL:   viper.GetDuration("jwt.token_ttl"),
		SigningKey: os.Getenv("SIGNING_KEY"),
	}

	token, err := service.NewAuthService(authConfig).GenerateToken(role)
	if err != nil {
		t.Fatalf("Error generate token: %s", err.Error())
	}

	return token
}

func SetAuthToken(r *http.Request, token string) {
	r.Header.Add("Authorization", "Bearer "+token)
}

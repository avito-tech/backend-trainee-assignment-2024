package utils

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	"gta2024/pkg/handler"
	"gta2024/pkg/repository"
	"gta2024/pkg/service"
)

func InitRouter(db *sqlx.DB) *gin.Engine {
	authConfig := service.AuthConfig{
		TokenTTL:   viper.GetDuration("jwt.token_ttl"),
		SigningKey: os.Getenv("SIGNING_KEY"),
	}

	repos := repository.NewRepository()
	services := service.NewService(db, repos, authConfig)
	handlers := handler.NewHandler(services)
	return handlers.InitRoutes()
}

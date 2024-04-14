package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"gta2024/pkg/handler"
	"gta2024/pkg/repository"
	"gta2024/pkg/service"
)

func InitRouter(db *sqlx.DB) *gin.Engine {
	repos := repository.NewRepository()
	services := service.NewService(db, repos)
	handlers := handler.NewHandler(services)
	return handlers.InitRoutes()
}

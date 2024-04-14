package handler

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gta2024/docs"
	"gta2024/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.New()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.POST("/auth/sign-in", h.signIn)
	r.GET("/user_banner", h.GetUserBanner)
	r.GET("/banner", h.GetBanners)
	r.POST("/banner", h.CreateBanner)
	r.PATCH("/banner/:id", h.UpdateBanner)
	r.DELETE("/banner/:id", h.DeleteBanner)

	return r
}

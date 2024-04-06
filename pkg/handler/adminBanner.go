package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"gta2024/pkg/models"
)

func (h *Handler) GetBanners(c *gin.Context) {
	c.JSON(http.StatusOK, models.BannerGetResp200{})
}

func (h *Handler) CreateBanner(c *gin.Context) {
	c.JSON(http.StatusCreated, models.BannerPostResp201{})
}

func (h *Handler) UpdateBanner(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (h *Handler) DeleteBanner(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

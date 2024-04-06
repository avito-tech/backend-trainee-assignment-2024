package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gta2024/pkg/models"
)

func (h *Handler) GetUserBanner(c *gin.Context) {
	c.JSON(http.StatusOK, models.UserBannerGetResp200{})
}

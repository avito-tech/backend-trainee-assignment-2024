package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"gta2024/pkg/models"
)

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Warn(message)
	c.AbortWithStatusJSON(statusCode, models.ErrorResp{Error: message})
}

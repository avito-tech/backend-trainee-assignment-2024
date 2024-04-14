package handler

import (
	"gta2024/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SignIn godoc
//
//	@Summary		Получение токена
//	@Description	Get token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			role	query		string						true	"Роль"	Enums(admin, user)
//	@Success		200		{object}	models.AuthSignInResp200	"OK"
//	@Failure		400		{object}	models.ErrorResp			"Bad Request"
//	@Failure		500		{object}	models.ErrorResp			"Internal Error"
//	@Router			/auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	role := c.Query("role")

	token, err := h.services.Authorization.GenerateToken(role)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, models.AuthSignInResp200{AccessToken: token})
}

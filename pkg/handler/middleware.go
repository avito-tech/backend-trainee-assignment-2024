package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userIsAdmin         = "isAdmin"
	adminRole           = "admin"
)

func (h *Handler) roleIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	role, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userIsAdmin, role == adminRole)
}

func (h *Handler) adminIdentity(c *gin.Context) {
	if !getIsAdmin(c) {
		newErrorResponse(c, http.StatusForbidden, "forbidden")
		return
	}
}

func getIsAdmin(c *gin.Context) bool {
	return c.GetBool(userIsAdmin)
}

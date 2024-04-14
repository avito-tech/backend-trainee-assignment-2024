package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gta2024/pkg/service"
)

// GetUserBanners godoc
//
//	@Security JWT
//	@Summary		Получение баннера для пользователя
//	@Description	Get user banners
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			tag_id				query		int64				true	"Тэг пользователя"
//	@Param			feature_id			query		int64				true	"Идентификатор фичи"
//	@Param			use_last_revision	query		bool				false	"Получать актуальную информацию"
//	@Success		200					{object}	models.ErrorResp	"OK"
//	@Failure		400					{object}	models.ErrorResp	"Bad Request"
//	@Failure		401					{object}	models.ErrorResp	"Unauthorized"
//	@Failure		403					{object}	models.ErrorResp	"Forbidden"
//	@Failure		404					{object}	models.ErrorResp	"Not Found"
//	@Failure		500					{object}	models.ErrorResp	"Internal Error"
//	@Router			/user_banner [get]
func (h *Handler) GetUserBanner(c *gin.Context) {
	var (
		useLastRevision = false
		isAdmin         = getIsAdmin(c)
		featureId       = parseQueryInt64Opt(c, "feature_id")
		tagId           = parseQueryInt64Opt(c, "tag_id")
	)

	if featureId == nil {
		newErrorResponse(c, http.StatusBadRequest, "feature_id invalid or missing")
		return
	}
	if tagId == nil {
		newErrorResponse(c, http.StatusBadRequest, "tag_id invalid or missing")
		return
	}

	if val := parseQueryBoolOpt(c, "use_last_revision"); val != nil {
		useLastRevision = *val
	}

	bannerContent, err := h.services.UserBanner.Get(*featureId, *tagId, isAdmin, useLastRevision)
	if err != nil {
		if errors.Is(err, service.ErrBannerNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, bannerContent)
}

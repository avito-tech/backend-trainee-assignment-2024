package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gta2024/pkg/models"
	"gta2024/pkg/service"
)

// GetBanners godoc
//
//	@Security JWT
//	@Summary		Получение всех баннеров c фильтрацией по фиче и/или тегу
//	@Description	Get banners
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			feature_id	query		int64					false	"Идентификатор фичи"
//	@Param			tag_id		query		int64					false	"Идентификатор тега"
//	@Param			limit		query		int64					false	"Лимит"
//	@Param			offset		query		int64					false	"Оффсет"
//	@Success		200			{object}	models.BannerGetResp200	"OK"
//	@Failure		401			{object}	models.ErrorResp		"Unauthorized"
//	@Failure		403			{object}	models.ErrorResp		"Forbidden"
//	@Failure		500			{object}	models.ErrorResp		"Internal Error"
//	@Router			/banner [get]
func (h *Handler) GetBanners(c *gin.Context) {
	feature_id := parseQueryInt64Opt(c, "feature_id")
	tag_id := parseQueryInt64Opt(c, "tag_id")
	limit := parseQueryInt64Opt(c, "limit")
	offset := parseQueryInt64Opt(c, "offset")

	banners, err := h.services.AdminBanner.GetAll(feature_id, tag_id, limit, offset)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, models.BannerGetResp200(banners))
}

// CreateBanner godoc
//
//	@Security JWT
//	@Summary		Создание нового баннера
//	@Description	Create banner
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CreateBanner			true	"CreateBanner"
//	@Success		201		{object}	models.BannerPostResp201	"Created"
//	@Failure		400		{object}	models.ErrorResp			"Bad Request"
//	@Failure		401		{object}	models.ErrorResp			"Unauthorized"
//	@Failure		403		{object}	models.ErrorResp			"Forbidden"
//	@Failure		409		{object}	models.ErrorResp			"Conflict"
//	@Failure		500		{object}	models.ErrorResp			"Internal Error"
//	@Router			/banner [post]
func (h *Handler) CreateBanner(c *gin.Context) {
	var createBanner models.CreateBanner
	if err := c.ShouldBindJSON(&createBanner); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	bannerId, err := h.services.AdminBanner.Create(createBanner)
	if err != nil {
		if errors.Is(err, service.ErrFeatureTagAlreadyExists) {
			newErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, models.BannerPostResp201{BannerID: bannerId})
}

// UpdateBanner godoc
//
//	@Security JWT
//	@Summary		Обновление содержимого баннера
//	@Description	Update banner
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int64				true	"Идентификатор баннера"
//	@Param			request	body		models.UpdateBanner	true	"UpdateBanner"
//	@Success		200		{string}	string				"Created"
//	@Failure		400		{object}	models.ErrorResp	"Bad Request"
//	@Failure		401		{object}	models.ErrorResp	"Unauthorized"
//	@Failure		403		{object}	models.ErrorResp	"Forbidden"
//	@Failure		404		{object}	models.ErrorResp	"Not Found"
//	@Failure		409		{object}	models.ErrorResp	"Conflict"
//	@Failure		500		{object}	models.ErrorResp	"Internal Error"
//	@Router			/banner/{id} [patch]
func (h *Handler) UpdateBanner(c *gin.Context) {
	bannerId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid banner id")
		return
	}
	var updateBanner models.UpdateBanner
	if err := c.ShouldBindJSON(&updateBanner); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if updateBanner.TagIDs != nil && len(*updateBanner.TagIDs) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "Empty tag_ids")
		return
	}
	err = h.services.AdminBanner.Update(bannerId, updateBanner)
	if err != nil {
		if errors.Is(err, service.ErrBannerNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, service.ErrFeatureTagAlreadyExists) {
			newErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// DeleteBanner godoc
//
//	@Security JWT
//	@Summary		Удаление баннера по идентификатору
//	@Description	Delete banner
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64				true	"Идентификатор баннера"
//	@Success		204	{string}	string				"Success Response"
//	@Failure		400	{object}	models.ErrorResp	"Bad Request"
//	@Failure		401	{object}	models.ErrorResp	"Unauthorized"
//	@Failure		403	{object}	models.ErrorResp	"Forbidden"
//	@Failure		404	{object}	models.ErrorResp	"Not Found"
//	@Failure		500	{object}	models.ErrorResp	"Internal Error"
//	@Router			/banner/{id} [delete]
func (h *Handler) DeleteBanner(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

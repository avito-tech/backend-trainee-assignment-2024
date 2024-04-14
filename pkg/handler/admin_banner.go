package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gta2024/pkg/models"
	"gta2024/pkg/service"
)

func parseQueryInt64Opt(c *gin.Context, name string) *int64 {
	valStr, ok := c.GetQuery(name)
	if !ok {
		return nil
	}
	val, _ := strconv.ParseInt(valStr, 10, 64)
	return &val
}

// GetBanners godoc
//
//		@Summary		Get banners
//		@Description	Get banners
//		@Tags			banner
//		@Accept			json
//		@Produce		json
//		@Param			request	 query		int64		false "feature_id"
//		@Param			request	 query		int64		false "tag_id"
//		@Param			request	 query		int64		false "limit"
//		@Param			request	 query		int64		false "offset"
//		@Success		200		{object}	models.ErrorResp "OK"
//	 @Failure        401 	{object}  	models.ErrorResp "Unauthorized"
//	 @Failure        403  	{object}  	models.ErrorResp "Forbidden"
//	 @Failure        403  	{object}  	models.ErrorResp "Internal Error"
//		@Router			/banner [get]
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

func (h *Handler) DeleteBanner(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

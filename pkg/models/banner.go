package models

import "time"

type (
	BannerContent map[string]interface{}

	Banner struct {
		BannerID  int           `json:"banner_id"`
		TagIDs    []int         `json:"tag_ids"`
		FeatureID int           `json:"feature_id"`
		Content   BannerContent `json:"content"`
		IsActive  bool          `json:"is_active"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
	}

	CreateBanner struct {
		TagIDs    []int         `json:"tag_ids"`
		FeatureID int           `json:"feature_id"`
		Content   BannerContent `json:"content"`
		IsActive  bool          `json:"is_active"`
	}

	UpdateBanner struct {
		TagIDs    *[]int         `json:"tag_ids"`
		FeatureID *int           `json:"feature_id"`
		Content   *BannerContent `json:"content"`
		IsActive  *bool          `json:"is_active"`
	}

	ErrorResp struct {
		Error string `json:"error"`
	}

	// GET /user_banner
	UserBannerGetResp200 BannerContent
	UserBannerGetResp400 ErrorResp
	UserBannerGetResp500 ErrorResp

	// GET /banner
	BannerGetResp200 []Banner
	BannerGetResp500 ErrorResp

	// POST /banner
	BannerPostReq     CreateBanner
	BannerPostResp201 struct {
		BannerID int `json:"banner_id"`
	}
	BannerPostResp400 ErrorResp
	BannerPostResp500 ErrorResp

	// PATCH /banner/{id}
	BannerIdPatchReq     UpdateBanner
	BannerIdPatchResp400 ErrorResp
	BannerIdPatchResp500 ErrorResp

	// DELETE /banner/{id}
	BannerIdDeleteResp400 ErrorResp
	BannerIdDeleteResp500 ErrorResp
)

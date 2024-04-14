package models

import "time"

type (
	DBBanner struct {
		BannerID  int64     `db:"id"`
		Content   string    `db:"content"`
		IsActive  bool      `db:"is_active"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	DBFeatureTagBanner struct {
		BannerID  int64 `db:"banner_id"`
		TagID     int64 `db:"tag_id"`
		FeatureID int64 `db:"feature_id"`
	}
)

type (
	BannerContent map[string]interface{}

	Banner struct {
		BannerID  int64         `json:"banner_id"`
		TagIDs    []int64       `json:"tag_ids"`
		FeatureID int64         `json:"feature_id"`
		Content   BannerContent `json:"content"`
		IsActive  bool          `json:"is_active"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
	}

	CreateBanner struct {
		TagIDs    *[]int64       `json:"tag_ids" binding:"required,min=1"`
		FeatureID *int64         `json:"feature_id" binding:"required"`
		Content   *BannerContent `json:"content" binding:"required"`
		IsActive  *bool          `json:"is_active" binding:"required"`
	}

	UpdateBanner struct {
		TagIDs    *[]int64       `json:"tag_ids"`
		FeatureID *int64         `json:"feature_id"`
		Content   *BannerContent `json:"content"`
		IsActive  *bool          `json:"is_active"`
	}

	ErrorResp struct {
		Error string `json:"error"`
	}

	// GET /user_banner
	UserBannerGetResp200 BannerContent

	// GET /banner
	BannerGetResp200 []Banner

	// POST /banner
	BannerPostReq     CreateBanner
	BannerPostResp201 struct {
		BannerID int64 `json:"banner_id"`
	}

	// PATCH /banner/{id}
	BannerIdPatchReq UpdateBanner
)

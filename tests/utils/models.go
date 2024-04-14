package utils

import (
	"github.com/jmoiron/sqlx"

	"gta2024/pkg/models"
)

func CreateOrUpdateBanner(db *sqlx.DB, banner models.DBBanner) error {
	query := `
		insert into banner (id, content, is_active, created_at, updated_at)
		values ($1, $2::jsonb, $3, $4, $5)
		on conflict (id) do update
			set content = $2,
				is_active = $3,
				updated_at = now()
	`
	_, err := db.Exec(query, banner.BannerID, banner.Content, banner.IsActive, banner.CreatedAt, banner.UpdatedAt)
	return err
}

func GetBanners(db *sqlx.DB) ([]models.DBBanner, error) {
	query := `
		select id, content, is_active, created_at, updated_at from banner
	`
	var banners []models.DBBanner
	err := db.Select(&banners, query)
	return banners, err
}

func CreateFeatureTagBanner(db *sqlx.DB, featureTagBanner models.DBFeatureTagBanner) error {
	query := `
		insert into feature_tag_banner (feature_id, tag_id, banner_id)
		values ($1, $2, $3)
		on conflict (feature_id, tag_id) do nothing
	`
	if _, err := db.Exec(query, featureTagBanner.FeatureID, featureTagBanner.TagID, featureTagBanner.BannerID); err != nil {
		return err
	}

	return nil
}

func GetFeatureTagBanners(db *sqlx.DB) ([]models.DBFeatureTagBanner, error) {
	query := `
		select feature_id, tag_id, banner_id from feature_tag_banner
	`
	var ftbs []models.DBFeatureTagBanner
	err := db.Select(&ftbs, query)
	return ftbs, err
}

package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"gta2024/pkg/models"
)

type FeatureTagBannerPostgres struct{}

func NewFeatureTagBannerPostgres() *FeatureTagBannerPostgres {
	return &FeatureTagBannerPostgres{}
}

func (r *FeatureTagBannerPostgres) GetBannerIdsByFeatureTag(db *sqlx.DB, featureId *int64, tagId *int64, limit int64, offset int64) ([]int64, error) {
	query := `
		select distinct banner_id from feature_tag_banner 
		where ($1::bigint is null or feature_id = $1)
		  and ($2::bigint is null or tag_id = $2)
		limit $3 offset $4
	`
	bannerIds := []int64{}
	err := db.Select(&bannerIds, query, featureId, tagId, limit, offset)
	if err != nil {
		return bannerIds, fmt.Errorf("error ftb.GetBannerIdsByFeatureTag: %w", err)
	}
	return bannerIds, nil
}

func (r *FeatureTagBannerPostgres) GetByBannerIds(db *sqlx.DB, bannerIds []int64) ([]models.DBFeatureTagBanner, error) {
	query := `
		select banner_id, tag_id, feature_id from feature_tag_banner 
		where banner_id in (select * from unnest($1::bigint[]))
	`
	dbFTBs := []models.DBFeatureTagBanner{}
	err := db.Select(&dbFTBs, query, pq.Array(bannerIds))
	if err != nil {
		return dbFTBs, fmt.Errorf("error ftb.GetByBannerIds: %w", err)
	}
	return dbFTBs, nil
}

func (r *FeatureTagBannerPostgres) Create(tx *sqlx.Tx, featureId int64, tagIds []int64, bannerId int64) error {
	query := `
		insert into feature_tag_banner (feature_id, tag_id, banner_id)
		select $1, tag_id, $3 from unnest($2::bigint[]) as src(tag_id)
	`
	_, err := tx.Exec(query, featureId, pq.Array(tagIds), bannerId)
	if err != nil {
		return fmt.Errorf("error ftb.Create: %w", err)
	}
	return err
}

func (r *FeatureTagBannerPostgres) GetByBannerId(tx *sqlx.Tx, bannerId int64) ([]models.DBFeatureTagBanner, error) {
	query := `select banner_id, tag_id, feature_id from feature_tag_banner 
		where banner_id = $1`
	dbFTBs := []models.DBFeatureTagBanner{}
	err := tx.Select(&dbFTBs, query, bannerId)
	if err != nil {
		return dbFTBs, fmt.Errorf("error ftb.GetByBannerId: %w", err)
	}
	return dbFTBs, nil
}

func (r *FeatureTagBannerPostgres) DeleteByBannerId(tx *sqlx.Tx, bannerId int64) error {
	query := `
		delete from feature_tag_banner 
		where banner_id = $1
	`
	_, err := tx.Exec(query, bannerId)
	if err != nil {
		return fmt.Errorf("error ftb.Delete: %w", err)
	}
	return nil
}

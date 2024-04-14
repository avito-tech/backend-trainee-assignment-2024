package repository

import (
	"github.com/jmoiron/sqlx"

	"gta2024/pkg/models"
)

type Banner interface {
	GetById(tx *sqlx.Tx, bannerId int64) (models.DBBanner, error)
	GetByIds(db *sqlx.DB, bannerIds []int64) ([]models.DBBanner, error)
	Create(tx *sqlx.Tx, content models.BannerContent, isActive bool) (int64, error)
	Update(tx *sqlx.Tx, banner models.DBBanner) error
	GetByFeatureIdTagId(db *sqlx.DB, featureId int64, tagId int64) (models.DBBanner, error)
}

type FeatureTagBanner interface {
	GetBannerIdsByFeatureTag(db *sqlx.DB, featureId, tagId *int64, limit int64, offset int64) ([]int64, error)
	GetByBannerIds(db *sqlx.DB, bannerIds []int64) ([]models.DBFeatureTagBanner, error)
	Create(tx *sqlx.Tx, featureId int64, tagIds []int64, bannerId int64) error
	GetByBannerId(tx *sqlx.Tx, bannerId int64) ([]models.DBFeatureTagBanner, error)
	DeleteByBannerId(tx *sqlx.Tx, bannerId int64) error
}

type Repository struct {
	Banner
	FeatureTagBanner
}

func NewRepository() *Repository {
	return &Repository{
		Banner:           NewBannerPostgres(),
		FeatureTagBanner: NewFeatureTagBannerPostgres(),
	}
}

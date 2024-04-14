package service

import (
	"github.com/jmoiron/sqlx"

	"gta2024/pkg/models"
	"gta2024/pkg/repository"
)

type UserBanner interface {
	Get(tagId int, featureId int, useLastRevision bool) (models.BannerContent, error)
}

type AdminBanner interface {
	Create(banner models.CreateBanner) (int64, error)
	GetAll(featureId, tagId, limit, offset *int64) ([]models.Banner, error)
	Update(bannerId int64, banner models.UpdateBanner) error
}

type Service struct {
	// UserBanner
	AdminBanner
}

func NewService(db *sqlx.DB, repos *repository.Repository) *Service {
	return &Service{
		// UserBanner:  NewAuthService(repos.Authorization),
		AdminBanner: NewAdminBannerService(db, repos.Banner, repos.FeatureTagBanner),
	}
}

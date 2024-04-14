package service

import (
	"github.com/jmoiron/sqlx"

	"gta2024/pkg/models"
	"gta2024/pkg/repository"
)

type UserBanner interface {
	Get(featureId int64, tagId int64, isAdmin bool, useLastRevision bool) (models.BannerContent, error)
}

type AdminBanner interface {
	Create(banner models.CreateBanner) (int64, error)
	GetAll(featureId, tagId, limit, offset *int64) ([]models.Banner, error)
	Update(bannerId int64, banner models.UpdateBanner) error
}

type Service struct {
	UserBanner
	AdminBanner
}

func NewService(db *sqlx.DB, repos *repository.Repository) *Service {
	return &Service{
		UserBanner:  NewUserBannerService(db, repos.Banner),
		AdminBanner: NewAdminBannerService(db, repos.Banner, repos.FeatureTagBanner),
	}
}

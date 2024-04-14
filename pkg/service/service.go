package service

import (
	"github.com/jmoiron/sqlx"

	"gta2024/pkg/models"
	"gta2024/pkg/repository"
)

type Authorization interface {
	GenerateToken(role string) (string, error)
	ParseToken(token string) (string, error)
}

type UserBanner interface {
	Get(featureId int64, tagId int64, isAdmin bool, useLastRevision bool) (models.BannerContent, error)
}

type AdminBanner interface {
	Create(banner models.CreateBanner) (int64, error)
	GetAll(featureId, tagId, limit, offset *int64) ([]models.Banner, error)
	Update(bannerId int64, banner models.UpdateBanner) error
	Delete(bannerId int64) error
}

type Service struct {
	Authorization
	UserBanner
	AdminBanner
}

func NewService(db *sqlx.DB, repos *repository.Repository, authConfig AuthConfig) *Service {
	return &Service{
		Authorization: NewAuthService(authConfig),
		UserBanner:    NewUserBannerService(db, repos.Banner),
		AdminBanner:   NewAdminBannerService(db, repos.Banner, repos.FeatureTagBanner),
	}
}

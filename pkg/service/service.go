package service

import "gta2024/pkg/models"

type UserBanner interface {
	Get(tagId int, featureId int, useLastRevision bool) (models.BannerContent, error)
}

type AdminBanner interface {
	GetAll(featureId, tagId, limit, offset *int) ([]models.Banner, error)
	Create(banner models.CreateBanner) (int, error)
	Update(bannerId int, banner models.UpdateBanner) error
	Delete(bannerId int) error
}

type Service struct {
	// UserBanner
	// AdminBanner
}

func NewService() *Service {
	return &Service{
		// UserBanner:  NewAuthService(repos.Authorization),
		// AdminBanner: NewTodoListService(repos.TodoList),
	}
}

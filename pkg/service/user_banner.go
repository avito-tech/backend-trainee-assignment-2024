package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"gta2024/pkg/models"
	"gta2024/pkg/repository"
)

type UserBannerService struct {
	db         *sqlx.DB
	bannerRepo repository.Banner
}

func NewUserBannerService(db *sqlx.DB, bannerRepo repository.Banner) *UserBannerService {
	return &UserBannerService{db: db, bannerRepo: bannerRepo}
}

func (s *UserBannerService) Get(featureId int64, tagId int64, isAdmin bool, useLastRevision bool) (models.BannerContent, error) {
	banner, err := s.bannerRepo.GetByFeatureIdTagId(s.db, featureId, tagId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBannerNotFound
		}
		return nil, fmt.Errorf("failed get user banner: %w", err)
	}

	if !banner.IsActive && !isAdmin {
		logrus.Info("Banner not active and user not admin")
		return nil, ErrBannerNotFound
	}

	var bannerContent models.BannerContent
	err = json.Unmarshal([]byte(banner.Content), &bannerContent)
	if err != nil {
		return nil, fmt.Errorf("error seriaize banner content: %w", err)
	}

	return bannerContent, nil
}

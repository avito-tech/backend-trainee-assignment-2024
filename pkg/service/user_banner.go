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
	db          *sqlx.DB
	bannerCache *Cache
	bannerRepo  repository.Banner
}

func NewUserBannerService(db *sqlx.DB, bannerRepo repository.Banner, bannerCache *Cache) *UserBannerService {
	return &UserBannerService{db: db, bannerRepo: bannerRepo, bannerCache: bannerCache}
}

func (s *UserBannerService) Get(featureId int64, tagId int64, isAdmin bool, useLastRevision bool) (models.BannerContent, error) {
	var banner models.DBBanner
	var isCacheMissing bool
	var err error

	if !useLastRevision {
		bannerRaw, ok := s.bannerCache.Get(s.GetBannerKey(featureId, tagId))
		if ok {
			banner = bannerRaw.(models.DBBanner)
		}
		isCacheMissing = !ok
	}

	if useLastRevision || isCacheMissing {
		logrus.Info("Cache missing or useLastRevision")

		banner, err = s.bannerRepo.GetByFeatureIdTagId(s.db, featureId, tagId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrBannerNotFound
			}
			return nil, fmt.Errorf("failed get user banner: %w", err)
		}

		s.bannerCache.Set(s.GetBannerKey(featureId, tagId), banner, 0)
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

func (s *UserBannerService) GetBannerKey(featureId int64, tagId int64) string {
	return fmt.Sprintf("%d::%d", featureId, tagId)
}

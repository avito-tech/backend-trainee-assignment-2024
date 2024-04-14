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

var defaultLimit = int64(10)
var defaultOffset = int64(0)

type AdminBannerService struct {
	db               *sqlx.DB
	bannerRepo       repository.Banner
	featureTagBanner repository.FeatureTagBanner
}

func NewAdminBannerService(db *sqlx.DB, bannerRepo repository.Banner, featureTagBanner repository.FeatureTagBanner) *AdminBannerService {
	return &AdminBannerService{db: db, bannerRepo: bannerRepo, featureTagBanner: featureTagBanner}
}

func (s *AdminBannerService) Create(banner models.CreateBanner) (int64, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %w", err)
	}

	var bannerId int64

	bannerId, err = s.bannerRepo.Create(tx, *banner.Content, *banner.IsActive)
	if err != nil {
		logrus.Errorf("Error create banner, rollback: %s", err.Error())

		if err := tx.Rollback(); err != nil {
			logrus.Errorf("Error rollback transaction: %s", err.Error())
			return 0, err
		}

		return 0, err
	}

	err = s.featureTagBanner.Create(tx, *banner.FeatureID, *banner.TagIDs, bannerId)
	if err != nil {
		logrus.Errorf("Error create feature tags banner, rollback: %s", err.Error())

		if err := tx.Rollback(); err != nil {
			logrus.Errorf("Error rollback transaction: %s", err.Error())
			return 0, err
		}

		return 0, ErrFeatureTagAlreadyExists
	}

	err = tx.Commit()
	if err != nil {
		logrus.Errorf("Error commit transaction: %s", err.Error())
		return 0, err
	}

	return bannerId, nil
}

func (s *AdminBannerService) GetAll(featureId, tagId, limit, offset *int64) ([]models.Banner, error) {
	var finalLimit int64 = defaultLimit
	var finalOffset int64 = defaultOffset
	if limit != nil {
		finalLimit = *limit
	}
	if offset != nil {
		finalOffset = *offset
	}
	bannerIds, err := s.featureTagBanner.GetBannerIdsByFeatureTag(s.db, featureId, tagId, finalLimit, finalOffset)
	if err != nil {
		return nil, fmt.Errorf("failed get bannerIds: %w", err)
	}

	dbBanners, err := s.bannerRepo.GetByIds(s.db, bannerIds)
	if err != nil {
		return nil, fmt.Errorf("failed get banners by ids: %w", err)
	}

	dbFTBs, err := s.featureTagBanner.GetByBannerIds(s.db, bannerIds)
	if err != nil {
		return nil, fmt.Errorf("failed get ftb by bannerIds: %w", err)
	}

	bannerIdToDbBanner := make(map[int64]*models.DBBanner)
	for _, dbBanner := range dbBanners {
		bannerIdToDbBanner[dbBanner.BannerID] = &dbBanner
	}

	bannerIdToTagIds := make(map[int64][]int64)
	bannerIdToFeatureId := make(map[int64]int64)
	for _, ftb := range dbFTBs {
		if _, ok := bannerIdToDbBanner[ftb.BannerID]; !ok {
			bannerIdToTagIds[ftb.BannerID] = []int64{}
		}

		bannerIdToTagIds[ftb.BannerID] = append(bannerIdToTagIds[ftb.BannerID], ftb.TagID)
		bannerIdToFeatureId[ftb.BannerID] = ftb.FeatureID
	}

	banners := make([]models.Banner, 0, len(bannerIds))
	for _, bannerId := range bannerIds {
		dbBanner := bannerIdToDbBanner[bannerId]
		tagIds := bannerIdToTagIds[bannerId]
		featureId := bannerIdToFeatureId[bannerId]

		var bannerContent models.BannerContent
		err = json.Unmarshal([]byte(dbBanner.Content), &bannerContent)
		if err != nil {
			logrus.Warnf("Error while deserialize content, skipping: %s", err.Error())
			continue
		}

		banner := models.Banner{
			BannerID:  dbBanner.BannerID,
			IsActive:  dbBanner.IsActive,
			CreatedAt: dbBanner.CreatedAt,
			UpdatedAt: dbBanner.UpdatedAt,
			Content:   bannerContent,
			TagIDs:    tagIds,
			FeatureID: featureId,
		}
		banners = append(banners, banner)
	}

	return banners, nil
}

func (s *AdminBannerService) Update(bannerId int64, banner models.UpdateBanner) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	dbBanner, err := s.bannerRepo.GetById(tx, bannerId)

	if err != nil {
		logrus.Warningf("Error get banner, rollback: %s", err.Error())

		if err := tx.Rollback(); err != nil {
			logrus.Errorf("Error rollback transaction: %s", err.Error())
			return err
		}
		if errors.Is(err, sql.ErrNoRows) {
			return ErrBannerNotFound
		}

		return err
	}

	if banner.Content != nil {
		bannerContent, err := json.Marshal(banner.Content)
		if err != nil {
			logrus.Warningf("Error serialize banner content, rollback: %s", err.Error())

			if err := tx.Rollback(); err != nil {
				logrus.Errorf("Error rollback transaction: %s", err.Error())
				return err
			}
			return err
		}
		dbBanner.Content = string(bannerContent)
	}

	if banner.IsActive != nil {
		dbBanner.IsActive = *banner.IsActive
	}

	err = s.bannerRepo.Update(tx, dbBanner)
	if err != nil {
		logrus.Errorf("Error update banner, rollback: %s", err.Error())

		if err := tx.Rollback(); err != nil {
			logrus.Errorf("Error rollback transaction: %s", err.Error())
			return err
		}

		return err
	}

	dbFTBs, err := s.featureTagBanner.GetByBannerId(tx, bannerId)
	if err != nil {
		logrus.Errorf("Error get feature tag banners, rollback: %s", err.Error())

		if err := tx.Rollback(); err != nil {
			logrus.Errorf("Error rollback transaction: %s", err.Error())
			return err
		}

		return err
	}

	if len(dbFTBs) == 0 {
		return errors.New("no rows in feature_tag_banner")
	}

	err = s.featureTagBanner.DeleteByBannerId(tx, bannerId)
	if err != nil {
		logrus.Errorf("Error delete feature tag banner, rollback: %s", err.Error())

		if err := tx.Rollback(); err != nil {
			logrus.Errorf("Error rollback transaction: %s", err.Error())
			return err
		}

		return err
	}

	var finalFeature int64 = dbFTBs[0].FeatureID
	if banner.FeatureID != nil {
		finalFeature = *banner.FeatureID
	}

	var finalTags []int64
	if banner.TagIDs != nil {
		finalTags = *banner.TagIDs
	} else {
		for _, dbFTB := range dbFTBs {
			finalTags = append(finalTags, dbFTB.TagID)
		}
	}

	err = s.featureTagBanner.Create(tx, finalFeature, finalTags, bannerId)
	if err != nil {
		logrus.Errorf("Error create feature tag banner, rollback: %s", err.Error())

		if err := tx.Rollback(); err != nil {
			logrus.Errorf("Error rollback transaction: %s", err.Error())
			return err
		}
		return ErrFeatureTagAlreadyExists
	}

	err = tx.Commit()
	if err != nil {
		logrus.Errorf("Error commit transaction: %s", err.Error())
		return err
	}

	return nil
}

func (s *AdminBannerService) Delete(bannerId int64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	_, err = s.bannerRepo.GetById(tx, bannerId)
	if err != nil {
		logrus.Warningf("Error get banner, rollback: %s", err.Error())

		if err := tx.Rollback(); err != nil {
			logrus.Errorf("Error rollback transaction: %s", err.Error())
			return err
		}
		if errors.Is(err, sql.ErrNoRows) {
			return ErrBannerNotFound
		}

		return err
	}

	err = s.featureTagBanner.DeleteByBannerId(tx, bannerId)
	if err != nil {
		logrus.Errorf("Error delete feature tag banner, rollback: %s", err.Error())

		if err := tx.Rollback(); err != nil {
			logrus.Errorf("Error rollback transaction: %s", err.Error())
			return err
		}

		return err
	}

	err = s.bannerRepo.Delete(tx, bannerId)
	if err != nil {
		logrus.Errorf("Error delete banner, rollback: %s", err.Error())

		if err := tx.Rollback(); err != nil {
			logrus.Errorf("Error rollback transaction: %s", err.Error())
			return err
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		logrus.Errorf("Error commit transaction: %s", err.Error())
		return err
	}

	return nil
}

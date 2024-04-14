package repository

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"gta2024/pkg/models"
)

type BannerPostgres struct{}

func NewBannerPostgres() *BannerPostgres {
	return &BannerPostgres{}
}

func (r *BannerPostgres) GetById(tx *sqlx.Tx, bannerId int64) (models.DBBanner, error) {
	query := `
		select id, content, is_active, created_at, updated_at from banner 
		where id = $1
	`
	banner := models.DBBanner{}
	err := tx.QueryRowx(query, bannerId).StructScan(&banner)
	if err != nil {
		return banner, fmt.Errorf("error b.GetById: %w", err)
	}
	return banner, err

}

func (r *BannerPostgres) GetByIds(db *sqlx.DB, bannerIds []int64) ([]models.DBBanner, error) {
	query := `
		select id, content, is_active, created_at, updated_at from banner 
		where id in (select * from unnest($1::bigint[]))
	`
	banners := []models.DBBanner{}
	err := db.Select(&banners, query, pq.Array(bannerIds))
	if err != nil {
		return banners, fmt.Errorf("error b.GetByIds: %w", err)
	}
	return banners, err
}

func (r *BannerPostgres) Create(tx *sqlx.Tx, content models.BannerContent, isActive bool) (int64, error) {
	var id int64

	contentBytes, err := json.Marshal(content)
	if err != nil {
		return 0, fmt.Errorf("error serialize banner content: %w", err)
	}

	query := `
		insert into banner (content, is_active)
		values ($1, $2)
		returning id
	`

	row := tx.QueryRow(query, contentBytes, isActive)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error insert banner: %w", err)
	}

	return id, nil
}

func (r *BannerPostgres) Update(tx *sqlx.Tx, banner models.DBBanner) error {
	query := `
		update banner
	 	set
			content = $2,
			is_active = $3,
			updated_at = now()
		where id = $1
	`
	_, err := tx.Exec(query, banner.BannerID, banner.Content, banner.IsActive)
	if err != nil {
		return fmt.Errorf("error update banner: %w", err)
	}
	return nil
}

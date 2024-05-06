package repos

import (
	"avito-backend-2024-trainee/internal/entity"
	"avito-backend-2024-trainee/pkg/postgresql"
	"context"
	"errors"
)

var (
	ErrBannerNotFound = errors.New("banner not found")
	BannerScanError   = errors.New("scanning to banner went wrong")
)

type BannerRepo struct {
	postgres postgresql.Postgresql
}

func (b *BannerRepo) GetUserBanner(ctx context.Context, tagId int, featureId int, useLatestVersion bool) (entity.BannerWithTag, error) {
	query := `
SELECT b.id, b.title, b.text, b.url, b.feature_id, bt.tag_id as tag
FROM banners b
         JOIN banner_tags bt ON b.id = bt.banner_id
WHERE bt.tag_id = $1 AND b.feature_id = $2;
`
	var result entity.BannerWithTag
	row := b.postgres.DB.Raw(query, tagId, featureId)
	row.Scan(&result)
	// TODO: переделать на нормальную проверку понимания того, нашли мы результат или нет
	if result.ID == 0 {
		return entity.BannerWithTag{}, ErrBannerNotFound
	}
	return result, nil
}

func NewBannerRepo(postgres postgresql.Postgresql) *BannerRepo {
	return &BannerRepo{
		postgres: postgres,
	}
}

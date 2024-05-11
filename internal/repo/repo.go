package repo

import (
	"avito-backend-2024-trainee/internal/entity"
	"avito-backend-2024-trainee/internal/repo/repos"
	"avito-backend-2024-trainee/internal/repo/repos/cache"
	"avito-backend-2024-trainee/pkg/postgresql"
	"context"
)

type Banner interface {
	GetUserBanner(ctx context.Context, tagId uint, featureId uint, useLastRevision bool) (entity.ProductionBanner, error)
	GetBanner(ctx context.Context, tagId, featureId, offset, limit *uint) ([]entity.BannerWithTags, error)
	CreateBanner(ctx context.Context, tagIds []uint, featureId uint, title, text, url string, isActive bool) (entity.BannerId, error)
	UpdateBanner(ctx context.Context, tagIds []uint, featureId uint, title, text, url string, isActive bool, bannerId uint) error
	DeleteBanner(ctx context.Context, bannerId uint) error
}
type Repositories struct {
	Banner
}

func NewRepositories(postgres postgresql.Postgresql, cache cache.BannerCache) *Repositories {
	return &Repositories{
		Banner: repos.NewBannerRepo(postgres, &cache),
	}
}

package repo

import (
	"avito-backend-2024-trainee/internal/entity"
	"avito-backend-2024-trainee/internal/repo/repos"
	"avito-backend-2024-trainee/pkg/postgresql"
	"context"
)

type Banner interface {
	GetUserBanner(ctx context.Context, tagId uint, featureId uint, useLatestVersion bool) (entity.BannerWithTag, error)
	GetBanner(ctx context.Context, tagId, featureId, offset, limit *uint) ([]entity.BannerWithTags, error)
	CreateBanner(ctx context.Context, tagIds []uint, featureId uint, title, text, url string, isActive bool) (entity.BannerId, error)
	UpdateBanner(ctx context.Context, tagIds []uint, featureId uint, title, text, url string, isActive bool, bannerId uint) error
	DeleteBanner(ctx context.Context, bannerId uint) error
}
type Repositories struct {
	Banner
}

func NewRepositories(postgres postgresql.Postgresql) *Repositories {
	return &Repositories{
		Banner: repos.NewBannerRepo(postgres),
	}
}

package repo

import (
	"avito-backend-2024-trainee/internal/entity"
	"avito-backend-2024-trainee/internal/repo/repos"
	"avito-backend-2024-trainee/pkg/postgresql"
	"context"
)

type Banner interface {
	GetUserBanner(ctx context.Context, tagId uint, featureId uint, useLatestVersion bool) (entity.BannerWithTag, error)
}
type Repositories struct {
	Banner
}

func NewRepositories(postgres postgresql.Postgresql) *Repositories {
	return &Repositories{
		Banner: repos.NewBannerRepo(postgres),
	}
}

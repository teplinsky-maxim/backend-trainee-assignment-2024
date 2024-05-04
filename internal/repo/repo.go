package repo

import (
	"avito-backend-2024-trainee/internal/repo/repos"
	"context"
)

type Banner interface {
	GetUserBanner(ctx context.Context, tagId int, featureId int, useLatestVersion bool)
}
type Repositories struct {
	Banner
}

func NewRepositories() *Repositories {
	return &Repositories{
		Banner: repos.NewBannerRepo(),
	}
}

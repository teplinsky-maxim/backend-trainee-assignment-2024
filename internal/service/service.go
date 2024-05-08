package service

import (
	"avito-backend-2024-trainee/internal/entity"
	"avito-backend-2024-trainee/internal/repo"
	"avito-backend-2024-trainee/internal/service/banner"
	"context"
)

type Services struct {
	Banner Banner
}
type Dependencies struct {
	Repositories repo.Repositories
}

type Banner interface {
	GetUserBanner(ctx context.Context, input *banner.GetUserBannerInput) (entity.BannerWithTag, error)
	GetBanner(ctx context.Context, input *banner.GetBannerInput) ([]entity.BannerWithTags, error)
	CreateBanner(ctx context.Context, input *banner.CreateBannerInput) (entity.BannerId, error)
	UpdateBanner(ctx context.Context, input *banner.UpdateBannerInput, bannerId uint) error
}

func NewServices(deps Dependencies) *Services {
	return &Services{
		Banner: banner.NewBannerService(deps.Repositories.Banner),
	}
}

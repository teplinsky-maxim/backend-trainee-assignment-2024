package banner

import (
	"avito-backend-2024-trainee/internal/entity"
	"avito-backend-2024-trainee/internal/repo"
	"context"
)

type BannerService struct {
	bannerRepo repo.Banner
}

func (b BannerService) GetUserBanner(ctx context.Context, input *GetUserBannerInput) (entity.BannerWithTag, error) {
	result, err := b.bannerRepo.GetUserBanner(ctx, input.TagId, input.FeatureId, input.UseLatestVersion)
	if err != nil {
		return entity.BannerWithTag{}, err
	}
	return result, nil
}

func (b BannerService) GetBanner(ctx context.Context, input *GetBannerInput) ([]entity.BannerWithTags, error) {
	result, err := b.bannerRepo.GetBanner(ctx, input.TagId, input.FeatureId, input.Limit, input.Offset)
	if err != nil {
		return []entity.BannerWithTags{}, err
	}
	return result, nil
}

func NewBannerService(bannerRepo repo.Banner) *BannerService {
	return &BannerService{
		bannerRepo: bannerRepo,
	}
}

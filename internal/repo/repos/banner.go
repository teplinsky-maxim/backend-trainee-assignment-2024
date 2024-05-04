package repos

import (
	"context"
)

type BannerRepo struct {
}

func (b BannerRepo) GetUserBanner(ctx context.Context, tagId int, featureId int, useLatestVersion bool) {
	//TODO implement me
	panic("implement me")
}

func NewBannerRepo() *BannerRepo {
	return &BannerRepo{}
}

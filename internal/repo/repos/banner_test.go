package repos

import (
	"avito-backend-2024-trainee/config"
	"avito-backend-2024-trainee/internal/entity"
	"avito-backend-2024-trainee/pkg/middleware/auth"
	"avito-backend-2024-trainee/pkg/postgresql"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserBanner(t *testing.T) {
	conf, err := config.NewConfigWithDiscover(nil)
	if err != nil {
		panic(err)
	}
	connection, err := postgresql.NewTestConnection(conf)
	if err != nil {
		panic(err)
	}
	sepCon := postgresql.Postgresql(connection)
	bannerRepo := NewBannerRepo(sepCon)
	featureId := uint(1)
	tagId := uint(1)
	banner := entity.Banner{
		Title:     "Test banner",
		Text:      "Test text",
		Url:       "https://123.com",
		FeatureId: featureId,
		IsActive:  true,
	}
	connection.DB.Create(&banner)
	connection.DB.Create(&entity.BannerTag{
		BannerId: int(banner.ID),
		TagId:    int(tagId),
	})
	ctx := context.WithValue(context.Background(), auth.RoleCtxField, auth.USER)
	userBanner, err := bannerRepo.GetUserBanner(ctx, int(tagId), int(featureId), false)
	if err != nil {
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, int(userBanner.Tag), int(tagId))
	assert.Equal(t, int(userBanner.FeatureId), int(featureId))
}

package repos

import (
	"avito-backend-2024-trainee/config"
	"avito-backend-2024-trainee/internal/entity"
	"avito-backend-2024-trainee/internal/utils/database"
	"avito-backend-2024-trainee/pkg/middleware/auth"
	"avito-backend-2024-trainee/pkg/postgresql"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserBanner(t *testing.T) {
	tests := []struct {
		name           string
		tagId          uint
		featureId      uint
		useLatestVer   bool
		expectedErr    error
		bannerIsActive bool
		userRole       auth.Role
	}{
		{
			name:           "valid banner",
			tagId:          uint(1),
			featureId:      uint(1),
			useLatestVer:   false,
			expectedErr:    nil,
			bannerIsActive: true,
			userRole:       auth.USER,
		},
		{
			name:           "banner not found 1",
			tagId:          uint(68768687),
			featureId:      uint(1),
			useLatestVer:   false,
			expectedErr:    ErrBannerNotFound,
			bannerIsActive: true,
			userRole:       auth.USER,
		},
		{
			name:           "banner not found 1",
			tagId:          uint(1),
			featureId:      uint(12398123),
			useLatestVer:   false,
			expectedErr:    ErrBannerNotFound,
			bannerIsActive: true,
			userRole:       auth.USER,
		},
		{
			name:           "banner is not active",
			tagId:          uint(2),
			featureId:      uint(2),
			useLatestVer:   false,
			expectedErr:    BannerIsNotActiveError,
			bannerIsActive: false,
			userRole:       auth.USER,
		},
		{
			name:           "admin can access not active banner",
			tagId:          uint(2),
			featureId:      uint(2),
			useLatestVer:   false,
			expectedErr:    BannerIsNotActiveError,
			bannerIsActive: false,
			userRole:       auth.ADMIN,
		},
	}

	conf, err := config.NewConfigWithDiscover(nil)
	if err != nil {
		panic(err)
	}
	connection, err := postgresql.NewTestConnection(conf)
	if err != nil {
		panic(err)
	}
	bannerRepo := NewBannerRepo(postgresql.Postgresql(connection))

	tableName, _ := database.GetTableName(postgresql.Postgresql(connection), entity.Banner{})
	connection.SetUp(tableName)
	tableName, _ = database.GetTableName(postgresql.Postgresql(connection), entity.BannerTag{})
	connection.SetUp(tableName)

	banner := entity.Banner{
		Title:     "Test banner",
		Text:      "Test text",
		Url:       "https://123.com",
		FeatureId: 1,
		IsActive:  true,
	}
	connection.DB.Create(&banner)
	connection.DB.Create(&entity.BannerTag{
		BannerId: banner.ID,
		TagId:    1,
	})

	banner2 := entity.Banner{
		Title:     "Test banner",
		Text:      "Test text",
		Url:       "https://123.com",
		FeatureId: 2,
		IsActive:  false,
	}
	connection.DB.Create(&banner2)
	connection.DB.Create(&entity.BannerTag{
		BannerId: banner2.ID,
		TagId:    2,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), auth.RoleCtxField, tt.userRole)
			userBanner, err := bannerRepo.GetUserBanner(ctx, tt.tagId, tt.featureId, tt.useLatestVer)
			if err != nil {
				assert.Equal(t, tt.expectedErr, err)
			} else {
				if tt.userRole == auth.ADMIN {
					assert.Equal(t, banner2.ID, userBanner.ID)
				} else {
					assert.Equal(t, banner.ID, userBanner.ID)
				}
			}
		})
	}
}

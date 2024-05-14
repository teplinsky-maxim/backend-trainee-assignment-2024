package repos

import (
	"avito-backend-2024-trainee/config"
	"avito-backend-2024-trainee/internal/entity"
	"avito-backend-2024-trainee/internal/repo/repos/cache"
	"avito-backend-2024-trainee/internal/utils/database"
	"avito-backend-2024-trainee/pkg/middleware/auth"
	"avito-backend-2024-trainee/pkg/postgresql"
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
	inMemoryCache := cache.NewInMemoryCache(5 * time.Minute)
	bannerRepo := NewBannerRepo(postgresql.Postgresql(connection), &inMemoryCache)

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
					assert.Equal(t, banner2.Text, userBanner.Text)
					assert.Equal(t, banner2.Title, userBanner.Title)
					assert.Equal(t, banner2.Url, userBanner.Url)
				} else {
					assert.Equal(t, banner.Text, userBanner.Text)
					assert.Equal(t, banner.Title, userBanner.Title)
					assert.Equal(t, banner.Url, userBanner.Url)
				}
			}
		})
	}
}

func BenchmarkGetUserBanner(b *testing.B) {
	conf, err := config.NewConfigWithDiscover(nil)
	if err != nil {
		panic(err)
	}
	connection, err := postgresql.NewTestConnection(conf)
	if err != nil {
		panic(err)
	}
	inMemoryCache := cache.NewInMemoryCache(5 * time.Minute)
	bannerRepo := NewBannerRepo(postgresql.Postgresql(connection), &inMemoryCache)

	// Parameters for GetUserBanner function
	tagId := uint(1)     // Replace with actual tag ID
	featureId := uint(1) // Replace with actual feature ID

	const (
		CommonTag1    = 5
		CommonTag2    = 10
		CommonFeature = 14
	)

	var table []struct{ tagId, featureId uint }
	for i := 0; i < 20; i++ {
		if i%4 == 0 {
			featureId = CommonFeature
		} else {
			featureId = gofakeit.UintRange(1, 1000)
		}

		if i%3 == 0 {
			tagId = CommonTag1
		} else if i%4 == 0 {
			tagId = CommonTag2
		} else {
			tagId = gofakeit.UintRange(1, 1000)
		}
		table = append(table, struct{ tagId, featureId uint }{
			tagId:     tagId,
			featureId: featureId,
		})
	}

	// Reset the benchmark timer
	b.ResetTimer()

	for idx, v := range table {
		b.Run(fmt.Sprintf("input_size_%d", idx), func(b *testing.B) {
			ctx := context.WithValue(context.Background(), auth.RoleCtxField, auth.USER)
			_, _ = bannerRepo.GetUserBanner(ctx, v.tagId, v.featureId, true)
			//featureId = gofakeit.UintRange(1, 1000)
			//title := gofakeit.BookTitle()
			//text := gofakeit.Sentence(gofakeit.IntRange(4, 10))
			//url := gofakeit.URL()
			//b.StopTimer()
			//bannerRepo.UpdateBanner(
			//	context.Background(), generateBannerTags(), featureId, title, text, url, true,
			//	gofakeit.UintRange(1, 10000),
			//)
			b.StartTimer()
		})
	}
}

func TestGetBanner(t *testing.T) {
	// разные параметры
	// проверить что все теги тянутся при where tagId =
}

func TestCreateBanner(t *testing.T) {
	//проверить разные значения
	//проверить, что тх роллбэчится на ошибке
}

func TestUpdateBanner(t *testing.T) {
	//проверить разные значения
	//упор на корнер-кейсы с тэгами
}

func TestDeleteBanner(t *testing.T) {
	//разные значения
}

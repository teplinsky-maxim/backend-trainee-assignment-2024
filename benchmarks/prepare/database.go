package prepare

import (
	"avito-backend-2024-trainee/config"
	"avito-backend-2024-trainee/internal/repo"
	"avito-backend-2024-trainee/internal/repo/repos/cache"
	"avito-backend-2024-trainee/pkg/postgresql"
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/samber/lo"
	"time"
)

const (
	CommonTag1    = 5
	CommonTag2    = 10
	CommonFeature = 14
)

func FillDatabase() {
	gofakeit.Seed(1111)
	conf, err := config.NewConfig(nil)
	if err != nil {
		panic(err)
	}
	postgres, err := postgresql.NewConnection(conf)
	if err != nil {
		panic(err)
	}
	inMemoryCache := cache.NewInMemoryCache(5 * time.Minute)
	repositories := repo.NewRepositories(postgres, &inMemoryCache)
	for i := 0; i < 10000; i++ {
		tags := generateBannerTags()
		var featureId uint
		if i%3 == 0 {
			featureId = CommonFeature
		} else {
			featureId = gofakeit.UintRange(1, 1000)
		}
		title := gofakeit.BookTitle()
		text := gofakeit.Sentence(gofakeit.IntRange(4, 10))
		url := gofakeit.URL()
		isActive := i%7 == 0
		_, err := repositories.Banner.CreateBanner(context.TODO(), tags, featureId, title, text, url, isActive)
		if err != nil {
			panic(err)
		}
	}
}

func generateBannerTags() []uint {

	amount := gofakeit.IntRange(4, 12)
	result := make([]uint, amount)
	for i := 0; i < amount; i++ {
		result[i] = gofakeit.UintRange(1, 1000)
	}
	result = append(result, CommonTag1)
	if gofakeit.UintRange(1, 10)%2 == 0 {
		result = append(result, CommonTag2)
	}
	result = lo.Shuffle(result)
	return result
}

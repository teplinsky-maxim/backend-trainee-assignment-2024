package app

import (
	"avito-backend-2024-trainee/config"
	v1 "avito-backend-2024-trainee/internal/controller/http/v1"
	"avito-backend-2024-trainee/internal/repo"
	"avito-backend-2024-trainee/internal/repo/repos/cache"
	"avito-backend-2024-trainee/internal/service"
	"avito-backend-2024-trainee/pkg/postgresql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

func Run() {
	//TODO: обернуть в singleton
	conf, err := config.NewConfig(nil)
	if err != nil {
		log.Fatal(err)
	}

	postgres, err := postgresql.NewConnection(conf)
	if err != nil {
		log.Fatal(err)
	}

	inMemoryCache := cache.NewInMemoryCache(5 * time.Minute)
	repositories := repo.NewRepositories(postgres, &inMemoryCache)

	serviceDeps := service.Dependencies{
		Repositories: *repositories,
	}
	services := service.NewServices(serviceDeps)

	app := fiber.New()
	v1.NewRouter(app, services)

	app.Listen("0.0.0.0:3000")
}

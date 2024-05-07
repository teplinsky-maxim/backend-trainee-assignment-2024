package app

import (
	"avito-backend-2024-trainee/config"
	v1 "avito-backend-2024-trainee/internal/controller/http/v1"
	"avito-backend-2024-trainee/internal/repo"
	"avito-backend-2024-trainee/internal/service"
	"avito-backend-2024-trainee/pkg/postgresql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Run(configPath string) {
	conf, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	postgres, err := postgresql.NewConnection(conf)
	if err != nil {
		log.Fatal(err)
	}
	//err = postgres.Migrate()
	//if err != nil {
	//	log.Fatal(err)
	//}

	repositories := repo.NewRepositories(postgres)

	deps := service.Dependencies{
		Repositories: *repositories,
	}
	services := service.NewServices(deps)

	app := fiber.New()
	v1.NewRouter(app, services)

	app.Listen("localhost:3000")
}

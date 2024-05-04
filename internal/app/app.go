package app

import (
	v1 "avito-backend-2024-trainee/internal/controller/http/v1"
	"avito-backend-2024-trainee/internal/service"
	"github.com/gofiber/fiber/v3"
)

func Run() {
	deps := service.ServiceDependencies{}
	services := service.NewServices(deps)

	app := fiber.New()
	v1.NewRouter(app, services)

	app.Listen("localhost:3000")
}

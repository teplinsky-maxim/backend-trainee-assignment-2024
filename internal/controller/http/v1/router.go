package v1

import (
	"avito-backend-2024-trainee/internal/service"
	"github.com/gofiber/fiber/v3"
)

func NewRouter(app *fiber.App, services *service.Services) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	NewBannerRoutes(&v1, services.Banner)
}

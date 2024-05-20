package v1

import (
	"avito-backend-2024-trainee/internal/service"
	"avito-backend-2024-trainee/pkg/middleware/auth"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, services *service.Services) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Use(auth.AuthMiddleware)
	NewBannerRoutes(&v1, services.Banner)
}

package v1

import (
	"avito-backend-2024-trainee/internal/service"
	"github.com/gofiber/fiber/v3"
)

type bannerRoutes struct {
	bannerService service.Banner
}

func NewBannerRoutes(router *fiber.Router, bannerService service.Banner) {
	r := &bannerRoutes{bannerService: bannerService}

	(*router).Add([]string{"GET"}, "/user_banner", r.create())
}

func (r *bannerRoutes) create() fiber.Handler {
	return func(c fiber.Ctx) error {
		r.bannerService.GetUserBanner()
		err := c.SendString("created")
		if err != nil {
			return err
		}
		return nil
	}
}

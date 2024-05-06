package v1

import (
	"avito-backend-2024-trainee/internal/repo/repos"
	"avito-backend-2024-trainee/internal/service"
	bannerService "avito-backend-2024-trainee/internal/service/banner"
	"avito-backend-2024-trainee/pkg/middleware/auth"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type bannerRoutes struct {
	bannerService service.Banner
}

func NewBannerRoutes(router *fiber.Router, bannerService service.Banner) {
	r := &bannerRoutes{bannerService: bannerService}

	(*router).Add("GET", "/user_banner/", r.create())
}

func (r *bannerRoutes) create() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: реализовать логику
		value := c.Locals(auth.ROLE_CTX_FIELD).(auth.Role)
		print(value)
		params := new(bannerService.GetUserBannerInput)
		if err := c.QueryParser(params); err != nil {
			return err
		}
		banner, err := r.bannerService.GetUserBanner(context.TODO(), params)
		if err != nil {
			if errors.Is(repos.ErrBannerNotFound, err) {
				_ = c.SendStatus(404)
				return nil
			} else if errors.Is(repos.BannerScanError, err) {
				_ = c.SendStatus(500)
				return nil
			}
			return err
		}
		err = c.JSON(banner)
		if err != nil {
			_ = c.SendStatus(500)
		}
		return nil
	}
}

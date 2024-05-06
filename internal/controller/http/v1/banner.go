package v1

import (
	"avito-backend-2024-trainee/internal/repo/repos"
	"avito-backend-2024-trainee/internal/service"
	banner2 "avito-backend-2024-trainee/internal/service/banner"
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
		params := new(banner2.GetUserBannerInput)
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

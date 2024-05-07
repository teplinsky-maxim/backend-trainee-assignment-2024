package v1

import (
	"avito-backend-2024-trainee/internal/repo/repos"
	"avito-backend-2024-trainee/internal/service"
	bannerService "avito-backend-2024-trainee/internal/service/banner"
	"avito-backend-2024-trainee/pkg/middleware/auth"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type bannerRoutes struct {
	bannerService service.Banner
}

func NewBannerRoutes(router *fiber.Router, bannerService service.Banner) {
	r := &bannerRoutes{bannerService: bannerService}

	(*router).Add("GET", "/user_banner/", r.getUserBannerHandler())
	(*router).Add("GET", "/banner/", r.getBannerHandler())
}

func (r *bannerRoutes) getUserBannerHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(bannerService.GetUserBannerInput)
		if err := c.QueryParser(params); err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		role := c.Locals(auth.RoleCtxField).(auth.Role)
		ctx := context.WithValue(context.Background(), auth.RoleCtxField, role)
		banner, err := r.bannerService.GetUserBanner(ctx, params)
		if err != nil {
			if errors.Is(repos.ErrBannerNotFound, err) {
				return c.SendStatus(http.StatusNotFound)
			} else if errors.Is(repos.BannerScanError, err) {
				return c.SendStatus(http.StatusInternalServerError)
			} else if errors.Is(repos.BannerIsNotActiveError, err) {
				return c.SendStatus(http.StatusForbidden)
			}
			return err
		}
		return c.JSON(banner)
	}
}

func (r *bannerRoutes) getBannerHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(bannerService.GetBannerInput)
		if err := c.QueryParser(params); err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		role := c.Locals(auth.RoleCtxField).(auth.Role)
		if role != auth.ADMIN {
			return c.SendStatus(http.StatusUnauthorized)
		}
		banners, err := r.bannerService.GetBanner(context.TODO(), params)
		if err != nil {
			return err
		}
		return c.JSON(banners)
	}
}

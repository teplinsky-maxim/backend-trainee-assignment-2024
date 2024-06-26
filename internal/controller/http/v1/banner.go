package v1

import (
	"avito-backend-2024-trainee/internal/repo/repos"
	"avito-backend-2024-trainee/internal/repo/repos/cache"
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
	cache         *cache.BannerCache
}

func NewBannerRoutes(router *fiber.Router, bannerService service.Banner) {
	r := &bannerRoutes{bannerService: bannerService}

	(*router).Add("GET", "/user_banner/", r.getUserBannerHandler())
	(*router).Add("GET", "/banner/", r.getBannerHandler())
	(*router).Add("POST", "/banner/", r.createBannerHandler())
	(*router).Add("PATCH", "/banner/:id", r.updateBannerHandler())
	(*router).Add("DELETE", "/banner/:id", r.deleteBannerHandler())
}

func sendError(c *fiber.Ctx, errorCode int, err error) error {
	return c.Status(errorCode).JSON(map[string]string{
		"error": err.Error(),
	})
}

func (r *bannerRoutes) getUserBannerHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := new(bannerService.GetUserBannerInput)
		if err := c.QueryParser(params); err != nil {
			return sendError(c, http.StatusBadRequest, err)
		}

		role := auth.GetRoleFromFiberCtx(c)
		ctx := context.WithValue(context.Background(), auth.RoleCtxField, role)
		banner, err := r.bannerService.GetUserBanner(ctx, params)
		if err != nil {
			if errors.Is(repos.ErrBannerNotFound, err) {
				return c.SendStatus(http.StatusNotFound)
			} else if errors.Is(repos.BannerScanError, err) {
				return sendError(c, http.StatusInternalServerError, err)
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

		role := auth.GetRoleFromFiberCtx(c)
		if !role.IsAdmin() {
			return c.SendStatus(http.StatusUnauthorized)
		}
		banners, err := r.bannerService.GetBanner(context.TODO(), params)
		if err != nil {
			return err
		}
		return c.JSON(banners)
	}
}

func (r *bannerRoutes) createBannerHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := new(bannerService.CreateBannerInput)
		if err := c.BodyParser(body); err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		role := auth.GetRoleFromFiberCtx(c)
		if !role.IsAdmin() {
			return c.SendStatus(http.StatusUnauthorized)
		}
		banner, err := r.bannerService.CreateBanner(context.TODO(), body)
		if err != nil {
			return err
		}
		return c.JSON(banner)
	}
}

func (r *bannerRoutes) updateBannerHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := new(bannerService.UpdateBannerInput)
		existingBannerId, err := c.ParamsInt("id", -1)
		if existingBannerId == -1 || err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}
		if err := c.BodyParser(body); err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		role := auth.GetRoleFromFiberCtx(c)
		if !role.IsAdmin() {
			return c.SendStatus(http.StatusUnauthorized)
		}
		err = r.bannerService.UpdateBanner(context.TODO(), body, uint(existingBannerId))
		if errors.Is(err, repos.ErrBannerNotFound) {
			return c.SendStatus(http.StatusNotFound)
		}
		return err
	}
}

func (r *bannerRoutes) deleteBannerHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		existingBannerId, err := c.ParamsInt("id", -1)
		if existingBannerId == -1 || err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}
		role := auth.GetRoleFromFiberCtx(c)
		if !role.IsAdmin() {
			return c.SendStatus(http.StatusUnauthorized)
		}
		err = r.bannerService.DeleteBanner(context.TODO(), &bannerService.DeleteBannerInput{}, uint(existingBannerId))
		if err == nil {
			return c.SendStatus(http.StatusNoContent)
		}
		if errors.Is(err, repos.ErrBannerNotFound) {
			return c.SendStatus(http.StatusNotFound)
		}
		return err
	}
}

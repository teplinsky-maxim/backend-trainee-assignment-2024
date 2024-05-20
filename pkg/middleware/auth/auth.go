package auth

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

type Role int32

const (
	ADMIN = iota
	USER
)

const (
	RoleCtxField = "user_role"
)

var (
	NotAuthenticatedError = errors.New("not authenticated")
	WrongTokenError       = errors.New("wrong token")
)

func AuthMiddleware(c *fiber.Ctx) error {
	headerValue := c.Get("Authorization", "")
	role, err := decideRole(headerValue)
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	setRoleToFiberCtx(c, role)
	return c.Next()
}

func decideRole(token string) (Role, error) {
	if token == "" {
		return -1, NotAuthenticatedError
	}
	role := strings.Split(token, " ")
	if len(role) != 2 {
		return -1, WrongTokenError
	}
	actualRole := role[1]

	if actualRole == "admin" {
		return ADMIN, nil
	} else if actualRole == "user" {
		return USER, nil
	}
	return -1, WrongTokenError
}

func (r Role) IsAdmin() bool {
	return r == ADMIN
}

func (r Role) IsUser() bool {
	return r == USER
}

func (r Role) Code() int {
	return int(r)
}

func setRoleToFiberCtx(c *fiber.Ctx, role Role) {
	c.Locals(RoleCtxField, role)
}

func GetRoleFromFiberCtx(c *fiber.Ctx) Role {
	role := c.Locals(RoleCtxField).(Role)
	return role
}

func GetRoleFromCtx(ctx *context.Context) Role {
	role := (*ctx).Value(RoleCtxField).(Role)
	return role
}

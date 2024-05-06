package auth

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type Role int32

const (
	ADMIN = iota
	USER
)

const (
	ROLE_CTX_FIELD = "user_role"
)

var (
	NotAuthenticatedError = errors.New("not authenticated")
	WrongTokenError       = errors.New("wrong token")
)

func AuthMiddleware(c *fiber.Ctx) error {
	headerValue := c.Get("Authorization", "")
	role, err := decideRole(headerValue)
	if err != nil {
		return c.SendStatus(401)
	}

	c.Locals(ROLE_CTX_FIELD, role)
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

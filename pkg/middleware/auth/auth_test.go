package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	app := fiber.New()
	testHandler := func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	}
	app.Use(AuthMiddleware)
	app.Get("/test", testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	req.Header.Set("Authorization", "Bearer")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	req.Header.Set("Authorization", "Bearer admin")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req.Header.Set("Authorization", "Bearer invalid")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestDecideRole(t *testing.T) {
	role, err := decideRole("")
	assert.Equal(t, Role(-1), role)
	assert.Equal(t, NotAuthenticatedError, err)

	role, err = decideRole("Bearer")
	assert.Equal(t, Role(-1), role)
	assert.Equal(t, WrongTokenError, err)

	role, err = decideRole("Bearer admin")
	assert.Equal(t, Role(ADMIN), role)
	assert.NoError(t, err)

	role, err = decideRole("Bearer user")
	assert.Equal(t, Role(USER), role)
	assert.NoError(t, err)

	role, err = decideRole("Bearer invalid")
	assert.Equal(t, Role(-1), role)
	assert.Equal(t, WrongTokenError, err)
}

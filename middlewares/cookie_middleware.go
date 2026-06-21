package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequireAuth(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token == "" {
		return c.Redirect("/login")
	}
	c.Locals("token", token)
	return c.Next()
}

func RedirectIfAuth(c *fiber.Ctx) error {
	if c.Cookies("token") != "" {
		return c.Redirect("/dashboard")
	}
	return c.Next()
}

func SetTokenCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
	})
}

func ClearTokenCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})
}

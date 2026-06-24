package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// RequireAuth dipakai untuk melindungi halaman web (/dashboard, /tasks, /badges).
// Kalau tidak ada cookie "token", user diarahkan ke halaman login.
func RequireAuth(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token == "" {
		return c.Redirect("/login")
	}
	c.Locals("token", token)
	return c.Next()
}

// RedirectIfAuth dipakai di halaman login dan register.
// Kalau user sudah punya cookie (sudah login), langsung redirect ke dashboard.
func RedirectIfAuth(c *fiber.Ctx) error {
	if c.Cookies("token") != "" {
		return c.Redirect("/dashboard")
	}
	return c.Next()
}

// SetTokenCookie menyimpan JWT ke cookie browser selama 24 jam.
// HTTPOnly: true supaya cookie tidak bisa dibaca JavaScript (lebih aman).
func SetTokenCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
	})
}

// ClearTokenCookie menghapus cookie token dengan cara set expires ke masa lalu.
func ClearTokenCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})
}

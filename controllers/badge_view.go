package controllers

import "github.com/gofiber/fiber/v2"

// BadgePage menampilkan halaman koleksi badge.
// Data badge dimuat secara async oleh app.js via /api/badges.
// Route: GET /badges
func BadgePage(c *fiber.Ctx) error {
	return c.Render("badges/index", fiber.Map{
		"Title": "Koleksi Badge",
		"Page":  "badges",
	}, "layouts/app")
}

package controllers

import "github.com/gofiber/fiber/v2"

// DashboardPage menampilkan halaman dashboard.
// Data statistik user dimuat secara async oleh app.js via /api/dashboard.
// Route: GET /dashboard
func DashboardPage(c *fiber.Ctx) error {
	return c.Render("dashboard/index", fiber.Map{
		"Title": "Dashboard",
		"Page":  "dashboard",
	}, "layouts/app")
}

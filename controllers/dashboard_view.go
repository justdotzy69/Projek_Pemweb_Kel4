package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func DashboardPage(c *fiber.Ctx) error {
	return c.Render("dashboard/index", fiber.Map{
		"Title": "Dashboard",
		"Page":  "dashboard",
	}, "layouts/app")
}

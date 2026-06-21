package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func TaskPage(c *fiber.Ctx) error {
	return c.Render("tasks/index", fiber.Map{
		"Title": "Papan Quest",
		"Page":  "tasks",
	}, "layouts/app")
}

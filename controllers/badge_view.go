package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func BadgePage(c *fiber.Ctx) error {
	return c.Render("badges/index", fiber.Map{
		"Title": "Koleksi Badge",
		"Page":  "badges",
	}, "layouts/app")
}

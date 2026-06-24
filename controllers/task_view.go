package controllers

import "github.com/gofiber/fiber/v2"

// TaskPage menampilkan halaman daftar tugas.
// Data tugas dimuat secara async oleh app.js via /api/tasks.
// Route: GET /tasks
func TaskPage(c *fiber.Ctx) error {
	return c.Render("tasks/index", fiber.Map{
		"Title": "Papan Quest",
		"Page":  "tasks",
	}, "layouts/app")
}

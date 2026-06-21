package routes

import (
	"Projek_Pemweb_Kel4/controllers"
	"Projek_Pemweb_Kel4/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Rute Publik
	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)

	// Rute Privat (Wajib pakai Token JWT)
	api.Get("/dashboard", middlewares.Protected(), controllers.GetDashboard)
	
	// Rute Kategori
	api.Post("/categories", middlewares.Protected(), controllers.CreateCategory)
	api.Get("/categories", middlewares.Protected(), controllers.GetCategories)

	// Rute Badge
	api.Post("/badges", middlewares.Protected(), controllers.CreateBadge)
	api.Get("/badges", middlewares.Protected(), controllers.GetBadges)

	// Rute Task (Pastikan dua baris ini ada)
	api.Post("/tasks", middlewares.Protected(), controllers.CreateTask)
	api.Get("/tasks", middlewares.Protected(), controllers.GetTasks)
	api.Put("/tasks/:id/complete", middlewares.Protected(), controllers.CompleteTask)
}
package routes

import (
	"Projek_Pemweb_Kel4/controllers"
	"Projek_Pemweb_Kel4/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes mendefinisikan semua endpoint API aplikasi
func SetupRoutes(app *fiber.App) {
	// Membuat grup utama dengan awalan "/api"
	api := app.Group("/api")

	// --------------------------------------------------
	// Rute Publik (Tanpa Token)
	// --------------------------------------------------
	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)

	// --------------------------------------------------
	// Rute Privat (Wajib pakai Token JWT)
	// --------------------------------------------------
	// Semua rute di bawah baris ini akan dilindungi oleh middleware Protected()
	api.Get("/dashboard", middlewares.Protected(), controllers.GetDashboard)
}
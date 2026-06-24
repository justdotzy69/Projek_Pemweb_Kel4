package routes

import (
	"Projek_Pemweb_Kel4/controllers"
	"Projek_Pemweb_Kel4/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	// ── Static files ──────────────────────────────────────────────────────────
	// Folder static/ bisa diakses browser di URL /static/
	app.Static("/static", "./static")

	// ── Halaman web (return HTML) ─────────────────────────────────────────────

	// Root redirect ke halaman login
	app.Get("/", func(c *fiber.Ctx) error { return c.Redirect("/login") })

	// Halaman auth — kalau sudah login langsung redirect ke /dashboard
	app.Get("/login",    middlewares.RedirectIfAuth, controllers.LoginPage)
	app.Post("/login",   controllers.LoginSubmit)
	app.Get("/register", middlewares.RedirectIfAuth, controllers.RegisterPage)
	app.Post("/register", controllers.RegisterSubmit)
	app.Get("/logout",   controllers.Logout)

	// Halaman utama — wajib login (dicek via cookie)
	app.Get("/dashboard", middlewares.RequireAuth, controllers.DashboardPage)
	app.Get("/tasks",     middlewares.RequireAuth, controllers.TaskPage)
	app.Get("/badges",    middlewares.RequireAuth, controllers.BadgePage)

	// ── API (return JSON) ─────────────────────────────────────────────────────
	api := app.Group("/api")

	// Auth — tidak butuh token
	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login",    controllers.Login)

	// Semua route di bawah ini butuh header: Authorization: Bearer <token>
	api.Get("/dashboard", middlewares.Protected(), controllers.GetDashboard)

	api.Post("/categories", middlewares.Protected(), controllers.CreateCategory)
	api.Get("/categories",  middlewares.Protected(), controllers.GetCategories)

	api.Post("/badges", middlewares.Protected(), controllers.CreateBadge)
	api.Get("/badges",  middlewares.Protected(), controllers.GetBadges)

	api.Post("/tasks",                middlewares.Protected(), controllers.CreateTask)
	api.Get("/tasks",                 middlewares.Protected(), controllers.GetTasks)
	api.Put("/tasks/:id/complete",    middlewares.Protected(), controllers.CompleteTask)
}

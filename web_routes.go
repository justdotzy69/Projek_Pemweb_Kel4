package main

import (
	"Projek_Pemweb_Kel4/controllers"
	"Projek_Pemweb_Kel4/middlewares"

	"github.com/gofiber/fiber/v2"
)

func setupWebRoutes(app *fiber.App) {
	app.Static("/static", "./static")

	app.Get("/", func(c *fiber.Ctx) error { return c.Redirect("/dashboard") })

	app.Get("/login",     middlewares.RedirectIfAuth, controllers.LoginPage)
	app.Post("/login",    middlewares.RedirectIfAuth, controllers.LoginSubmit)
	app.Get("/register",  middlewares.RedirectIfAuth, controllers.RegisterPage)
	app.Post("/register", middlewares.RedirectIfAuth, controllers.RegisterSubmit)
	app.Get("/logout",    controllers.Logout)

	web := app.Group("/", middlewares.RequireAuth)
	web.Get("/dashboard", controllers.DashboardPage)
	web.Get("/tasks",     controllers.TaskPage)
	web.Get("/badges",    controllers.BadgePage)
}

package main

import (
	"log"
	"os"

	"Projek_Pemweb_Kel4/database"
	"Projek_Pemweb_Kel4/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan nilai default.")
	}

	// 2. Koneksi database
	database.ConnectDB()

	// 3. Setup template engine untuk halaman web
	engine := html.New("./views", ".html")
	engine.AddFunc("add", func(a, b int) int { return a + b })
	engine.AddFunc("xpOf", func(diff string) int {
		switch diff {
		case "easy":   return 10
		case "medium": return 20
		case "hard":   return 30
		}
		return 0
	})
	engine.Reload(true)

	// 4. Inisialisasi Fiber dengan template engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// 5. Route API backend (tidak diubah sama sekali)
	routes.SetupRoutes(app)

	// 6. Route web frontend (ditambahkan di web_routes.go)
	setupWebRoutes(app)

	// 7. Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server berjalan di http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}

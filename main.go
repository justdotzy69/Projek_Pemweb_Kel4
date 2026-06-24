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
	// Load konfigurasi dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan nilai default.")
	}

	// Koneksi ke database MySQL dan jalankan auto-migrate
	database.ConnectDB()

	// Setup template engine: folder views/, ekstensi .html
	engine := html.New("./views", ".html")

	// Fungsi tambahan yang bisa dipanggil dari dalam template HTML
	engine.AddFunc("add", func(a, b int) int { return a + b })
	engine.AddFunc("xpOf", func(diff string) int {
		switch diff {
		case "easy":   return 10
		case "medium": return 20
		case "hard":   return 30
		}
		return 0
	})

	// Aktifkan reload otomatis template saat development
	engine.Reload(true)

	// Inisialisasi Fiber dengan template engine di atas
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Daftarkan semua route (web + API) ke app
	routes.SetupRoutes(app)

	// Baca port dari .env, default ke 3000 kalau tidak ada
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server berjalan di http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}

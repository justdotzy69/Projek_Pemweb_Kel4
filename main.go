package main

import (
	"log"
	"os"

	"Projek_Pemweb_Kel4/database"
	"Projek_Pemweb_Kel4/routes" // Tambahkan import package routes kita

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load variabel environment (Untuk mengambil nilai PORT)
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan nilai default.")
	}

	// 2. Jalankan fungsi koneksi database yang sudah kita buat
	database.ConnectDB()

	// 3. Inisialisasi aplikasi Fiber
	app := fiber.New()

	// 4. Buat satu route sederhana untuk pengetesan awal
	routes.SetupRoutes(app)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Selamat datang di API To-Do List Gamification Kelompok 4!",
		})
	})

	// 5. Tentukan Port (Ambil dari .env, atau gunakan 3000 jika kosong)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// 6. Jalankan server Fiber
	log.Printf("Menjalankan server di port %s...", port)
	err = app.Listen(":" + port)
	if err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}
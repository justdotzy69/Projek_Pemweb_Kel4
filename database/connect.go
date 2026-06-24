package database

import (
	"fmt"
	"log"
	"os"

	"Projek_Pemweb_Kel4/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB adalah koneksi global yang dipakai semua controller
var DB *gorm.DB

func ConnectDB() {
	// Load .env (kalau belum di-load di main.go, ini jadi fallback)
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan!")
	}

	// Baca konfigurasi database dari environment variables
	dbUser     := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost     := os.Getenv("DB_HOST")
	dbPort     := os.Getenv("DB_PORT")
	dbName     := os.Getenv("DB_NAME")

	// Susun DSN (Data Source Name) — format standar koneksi MySQL
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	// Buka koneksi via GORM
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database! Error: ", err)
	}

	fmt.Println("Koneksi ke database MySQL berhasil!")
	DB = database

	// Auto-migrate: GORM akan buat/update tabel sesuai struct di models/
	// Urutan penting — User dan Category harus ada sebelum Task (foreign key)
	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Badge{},
		&models.Task{},
	)
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database: ", err)
	}

	fmt.Println("Migrasi database selesai!")
}

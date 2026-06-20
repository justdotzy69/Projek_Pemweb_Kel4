package database

import (
	"fmt"
	"log"
	"os"

	"Projek_Pemweb_Kel4/models" // Tambahkan import models kita
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB adalah variabel global untuk menampung *instance* koneksi database
var DB *gorm.DB

// ConnectDB adalah fungsi untuk membuka koneksi ke MySQL
func ConnectDB() {
	// 1. Load konfigurasi dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan!")
	}

	// 2. Ambil nilai dari file .env
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 3. Susun string DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	// 4. Buka koneksi menggunakan GORM
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database! Error: ", err)
	}

	fmt.Println("Koneksi ke database MySQL berhasil!")

	// 5. Simpan koneksi ke variabel global DB
	DB = database

	// 6. Jalankan Auto-Migrate untuk membuat tabel
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

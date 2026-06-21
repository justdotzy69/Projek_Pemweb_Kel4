package controllers

import (
	"Projek_Pemweb_Kel4/database"
	"Projek_Pemweb_Kel4/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Struktur data khusus untuk menerima input dari user
type AuthInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ---------------------------------------------------------
// Fungsi untuk Registrasi User
// ---------------------------------------------------------
func Register(c *fiber.Ctx) error {
	var input AuthInput

	// 1. Membaca data JSON dari body request
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	// 2. Validasi input sederhana
	if input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Email dan password wajib diisi",
		})
	}

	// 3. Hash Password (enkripsi)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengenkripsi password",
		})
	}

	// 4. Susun data user baru
	newUser := models.User{
		Email:    input.Email,
		Password: string(hashPassword),
		Role:     "user", // Role default adalah 'user'
	}

	// 5. Simpan ke database
	result := database.DB.Create(&newUser)
	if result.Error != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Email sudah terdaftar atau terjadi kesalahan database",
		})
	}

	// 6. Kembalikan response sukses
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Registrasi berhasil",
		"data": fiber.Map{
			"id":    newUser.ID,
			"email": newUser.Email,
			"role":  newUser.Role,
		},
	})
}

// ---------------------------------------------------------
// Fungsi untuk Login User dan generate JWT
// ---------------------------------------------------------
func Login(c *fiber.Ctx) error {
	var input AuthInput

	// 1. Membaca data JSON dari body request
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	// 2. Cari user di database berdasarkan Email
	var user models.User
	result := database.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Email atau password salah",
		})
	}

	// 3. Bandingkan password input dengan hash di database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Email atau password salah",
		})
	}

	// 4. Generate JWT Token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal membuat token autentikasi",
		})
	}

	// 5. Kembalikan response sukses beserta Token JWT
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Login berhasil",
		"token":   tokenString,
	})
}
package controllers

import (
	"Projek_Pemweb_Kel4/database"
	"Projek_Pemweb_Kel4/models"

	"github.com/gofiber/fiber/v2"
)

// CreateBadge untuk menambahkan master badge baru
func CreateBadge(c *fiber.Ctx) error {
	var badge models.Badge

	// 1. Baca input JSON dari client
	if err := c.BodyParser(&badge); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	// 2. Validasi input sederhana
	if badge.Name == "" || badge.ImageURL == "" || badge.RequiredLevel <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Nama, Image URL, dan Required Level (minimal 1) wajib diisi",
		})
	}

	// 3. Simpan ke database
	result := database.DB.Create(&badge)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan badge (Mungkin nama sudah ada)",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Badge berhasil dibuat",
		"data":    badge,
	})
}

// GetBadges untuk mengambil semua daftar badge yang tersedia
func GetBadges(c *fiber.Ctx) error {
	var badges []models.Badge

	database.DB.Find(&badges)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Data badge berhasil diambil",
		"data":    badges,
	})
}
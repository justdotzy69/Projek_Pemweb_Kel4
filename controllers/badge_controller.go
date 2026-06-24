package controllers

import (
	"Projek_Pemweb_Kel4/database"
	"Projek_Pemweb_Kel4/models"

	"github.com/gofiber/fiber/v2"
)

// CreateBadge membuat data master badge baru di sistem.
// Badge ini nantinya bisa diperoleh user saat naik level.
// Route: POST /api/badges
func CreateBadge(c *fiber.Ctx) error {
	var badge models.Badge

	if err := c.BodyParser(&badge); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	if badge.Name == "" || badge.ImageURL == "" || badge.RequiredLevel <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Nama, Image URL, dan Required Level (minimal 1) wajib diisi",
		})
	}

	if err := database.DB.Create(&badge).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan badge (mungkin nama sudah ada)",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Badge berhasil dibuat",
		"data":    badge,
	})
}

// GetBadges mengambil semua badge yang tersedia di sistem.
// Route: GET /api/badges
func GetBadges(c *fiber.Ctx) error {
	var badges []models.Badge
	database.DB.Find(&badges)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Data badge berhasil diambil",
		"data":    badges,
	})
}

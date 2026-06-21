package controllers

import (
	"Projek_Pemweb_Kel4/database"
	"Projek_Pemweb_Kel4/models"

	"github.com/gofiber/fiber/v2"
)

// GetDashboard mengambil statistik user yang sedang login
func GetDashboard(c *fiber.Ctx) error {
	// 1. Ambil user_id dari memori lokal yang diset oleh middleware
	userID := c.Locals("user_id").(uint)

	var user models.User

	// 2. Cari data user di database, sekalian bawa data relasi Badges-nya
	// Kita gunakan .Select untuk menyembunyikan password agar lebih aman
	result := database.DB.Preload("Badges").Select("id, email, role, total_xp, current_level").First(&user, userID)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Data pengguna tidak ditemukan",
		})
	}

	// 3. Kembalikan data ke klien
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Data dashboard berhasil diambil",
		"data":    user,
	})
}
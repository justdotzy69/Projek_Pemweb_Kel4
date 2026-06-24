package controllers

import (
	"Projek_Pemweb_Kel4/database"
	"Projek_Pemweb_Kel4/models"

	"github.com/gofiber/fiber/v2"
)

// GetDashboard mengambil data statistik user yang sedang login.
// Password tidak ikut diambil karena pakai .Select() yang eksplisit.
// Route: GET /api/dashboard
func GetDashboard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var user models.User

	// Preload("Badges") = ikut ambil data badge yang dimiliki user
	// Select() = hanya ambil kolom yang diperlukan, password tidak ikut
	result := database.DB.
		Preload("Badges").
		Select("id, email, role, total_xp, current_level").
		First(&user, userID)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Data pengguna tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Data dashboard berhasil diambil",
		"data":    user,
	})
}

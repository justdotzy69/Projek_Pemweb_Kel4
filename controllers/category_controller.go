package controllers

import (
	"Projek_Pemweb_Kel4/database"
	"Projek_Pemweb_Kel4/models"

	"github.com/gofiber/fiber/v2"
)

// CreateCategory membuat kategori baru untuk mengelompokkan tugas.
// Route: POST /api/categories
func CreateCategory(c *fiber.Ctx) error {
	var category models.Category

	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	if category.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Nama kategori wajib diisi",
		})
	}

	if err := database.DB.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan kategori (mungkin nama sudah ada)",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Kategori berhasil dibuat",
		"data":    category,
	})
}

// GetCategories mengambil semua kategori yang tersedia.
// Route: GET /api/categories
func GetCategories(c *fiber.Ctx) error {
	var categories []models.Category
	database.DB.Find(&categories)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Data kategori berhasil diambil",
		"data":    categories,
	})
}

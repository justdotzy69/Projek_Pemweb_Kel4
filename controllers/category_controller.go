package controllers

import (
	"Projek_Pemweb_Kel4/database"
	"Projek_Pemweb_Kel4/models"

	"github.com/gofiber/fiber/v2"
)

// CreateCategory untuk menambahkan kategori baru
func CreateCategory(c *fiber.Ctx) error {
	var category models.Category

	// 1. Baca input dari user (berupa JSON: {"name": "Kuliah"})
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	// 2. Validasi input
	if category.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Nama kategori wajib diisi",
		})
	}

	// 3. Simpan ke database
	result := database.DB.Create(&category)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan kategori (Mungkin nama sudah ada)",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Kategori berhasil dibuat",
		"data":    category,
	})
}

// GetCategories untuk mengambil semua daftar kategori
func GetCategories(c *fiber.Ctx) error {
	var categories []models.Category

	// Cari semua data di tabel categories
	database.DB.Find(&categories)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Data kategori berhasil diambil",
		"data":    categories,
	})
}
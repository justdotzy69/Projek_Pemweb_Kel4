package controllers

import (
	"Projek_Pemweb_Kel4/database"
	"Projek_Pemweb_Kel4/helpers" // Tambahkan import helper
	"Projek_Pemweb_Kel4/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateTask untuk menambahkan tugas baru
func CreateTask(c *fiber.Ctx) error {
	// Ambil user_id dari middleware (JWT)
	userID := c.Locals("user_id").(uint)

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	// Set kepemilikan tugas dan status default
	task.UserID = userID
	task.Status = "pending"

	// Validasi input
	if task.Title == "" || task.Difficulty == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Judul dan tingkat kesulitan (easy/medium/hard) wajib diisi",
		})
	}

	if err := database.DB.Create(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal membuat tugas",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Tugas berhasil dibuat",
		"data":    task,
	})
}

// GetTasks untuk mengambil daftar tugas milik user yang login
func GetTasks(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var tasks []models.Task
	
	// Cari tugas milik user_id ini, dan tarik juga data Kategorinya (jika ada)
	database.DB.Preload("Category").Where("user_id = ?", userID).Find(&tasks)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Daftar tugas berhasil diambil",
		"data":    tasks,
	})
}

// CompleteTask untuk mengubah status tugas menjadi selesai dan memberikan reward
func CompleteTask(c *fiber.Ctx) error {
	taskID := c.Params("id") // Mengambil ID tugas dari URL (contoh: /api/tasks/1/complete)
	userID := c.Locals("user_id").(uint)

	// Mulai Database Transaction (Agar jika error di tengah jalan, data tidak terpotong)
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var task models.Task

		// 1. Cari tugas berdasarkan ID dan UserID
		if err := tx.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Tugas tidak ditemukan"})
		}

		// 2. Cek apakah tugas sudah selesai sebelumnya
		if task.Status == "completed" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Tugas sudah diselesaikan sebelumnya"})
		}

		// 3. Ubah status tugas menjadi selesai
		task.Status = "completed"
		if err := tx.Save(&task).Error; err != nil {
			return err
		}

		// 4. Proses Gamifikasi: Tambah XP dan Cek Level
		var user models.User
		if err := tx.Preload("Badges").First(&user, userID).Error; err != nil {
			return err
		}

		// Hitung XP yang didapat
		gainedXP := helpers.CalculateXP(task.Difficulty)
		user.TotalXP += gainedXP

		// Hitung Level baru
		newLevel := helpers.CalculateLevel(user.TotalXP)
		levelUp := false
		var newBadges []models.Badge

		if newLevel > user.CurrentLevel {
			user.CurrentLevel = newLevel
			levelUp = true

			// 5. Jika naik level, cari Badge yang syarat levelnya sudah terpenuhi
			tx.Where("required_level <= ?", newLevel).Find(&newBadges)

			// Masukkan badge baru ke user jika belum punya
			for _, badge := range newBadges {
				sudahPunya := false
				for _, userBadge := range user.Badges {
					if badge.ID == userBadge.ID {
						sudahPunya = true
						break
					}
				}
				if !sudahPunya {
					// Berikan badge ke user (insert ke tabel relasi many2many)
					tx.Model(&user).Association("Badges").Append(&badge)
				}
			}
		}

		// 6. Simpan update data user
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// Response Sukses
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Tugas berhasil diselesaikan!",
			"reward": fiber.Map{
				"gained_xp": gainedXP,
				"total_xp":  user.TotalXP,
				"level":     user.CurrentLevel,
				"level_up":  levelUp,
			},
		})
	})

	// Jika ada error di dalam transaction
	if err != nil {
		// Pengecekan agar tidak double response jika error sudah di-handle di dalam transaksi
		return err
	}

	return nil
}
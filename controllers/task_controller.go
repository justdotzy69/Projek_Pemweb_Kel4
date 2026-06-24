package controllers

import (
	"Projek_Pemweb_Kel4/database"
	"Projek_Pemweb_Kel4/helpers"
	"Projek_Pemweb_Kel4/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateTask membuat tugas baru milik user yang sedang login.
// Route: POST /api/tasks
func CreateTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	// Paksa ownership dan status awal — tidak boleh dikirim dari client
	task.UserID = userID
	task.Status = "pending"

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

// GetTasks mengambil semua tugas milik user yang sedang login.
// Data kategori ikut di-load sekaligus (Preload) agar tidak N+1 query.
// Route: GET /api/tasks
func GetTasks(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var tasks []models.Task
	database.DB.Preload("Category").Where("user_id = ?", userID).Find(&tasks)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Daftar tugas berhasil diambil",
		"data":    tasks,
	})
}

// CompleteTask menandai tugas sebagai selesai dan memproses reward gamifikasi.
// Seluruh proses dibungkus dalam satu database transaction agar atomik —
// kalau ada yang gagal di tengah jalan, semua perubahan dibatalkan.
// Route: PUT /api/tasks/:id/complete
func CompleteTask(c *fiber.Ctx) error {
	taskID := c.Params("id")
	userID := c.Locals("user_id").(uint)

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var task models.Task

		// 1. Pastikan tugas ada dan memang milik user ini
		if err := tx.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Tugas tidak ditemukan",
			})
		}

		// 2. Cegah complete ganda
		if task.Status == "completed" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Tugas sudah diselesaikan sebelumnya",
			})
		}

		// 3. Update status tugas
		task.Status = "completed"
		if err := tx.Save(&task).Error; err != nil {
			return err
		}

		// 4. Ambil data user beserta badge yang sudah dimiliki
		var user models.User
		if err := tx.Preload("Badges").First(&user, userID).Error; err != nil {
			return err
		}

		// 5. Tambah XP dan hitung level baru
		gainedXP := helpers.CalculateXP(task.Difficulty)
		user.TotalXP += gainedXP
		newLevel  := helpers.CalculateLevel(user.TotalXP)

		levelUp    := false
		var newBadges []models.Badge

		if newLevel > user.CurrentLevel {
			user.CurrentLevel = newLevel
			levelUp = true

			// 6. Cari badge yang syarat levelnya sudah terpenuhi
			tx.Where("required_level <= ?", newLevel).Find(&newBadges)

			// Berikan badge yang belum dimiliki user
			for _, badge := range newBadges {
				sudahPunya := false
				for _, userBadge := range user.Badges {
					if badge.ID == userBadge.ID {
						sudahPunya = true
						break
					}
				}
				if !sudahPunya {
					tx.Model(&user).Association("Badges").Append(&badge)
				}
			}
		}

		// 7. Simpan perubahan XP dan level user
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

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

	if err != nil {
		return err
	}
	return nil
}

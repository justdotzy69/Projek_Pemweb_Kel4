package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected adalah middleware untuk memvalidasi JWT
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil token dari header Authorization
		authHeader := c.Get("Authorization")

		// Jika header kosong, tolak akses
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Akses ditolak, token tidak ditemukan",
			})
		}

		// 2. Pastikan formatnya adalah "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Format token tidak valid",
			})
		}

		tokenString := parts[1]

		// 3. Parse dan validasi keaslian token menggunakan JWT_SECRET
		secretKey := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma enkripsinya sesuai
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secretKey), nil
		})

		// Jika error atau token palsu/kedaluwarsa
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		// 4. Ekstrak data (Klaim) dari dalam token
		claims := token.Claims.(jwt.MapClaims)

		// Simpan user_id ke Locals agar bisa diakses oleh Controller nanti
		// Nilai dari JWT berbentuk float64, kita convert jadi uint agar cocok dengan model kita
		c.Locals("user_id", uint(claims["user_id"].(float64)))

		// 5. Lolos pemeriksaan, lanjut ke fungsi Controller
		return c.Next()
	}
}
package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected adalah middleware untuk route API (/api/*).
// Cara kerja: baca header Authorization → validasi JWT → simpan user_id ke Locals.
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Akses ditolak, token tidak ditemukan",
			})
		}

		// 2. Format harus "Bearer <token>" — pisah jadi dua bagian
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Format token tidak valid",
			})
		}

		tokenString := parts[1]

		// 3. Parse dan verifikasi token menggunakan JWT_SECRET dari .env
		secretKey := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma yang dipakai adalah HMAC (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		// 4. Simpan user_id ke c.Locals supaya bisa diakses controller
		// Nilai dari JWT selalu float64, di-convert ke uint agar cocok dengan model
		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", uint(claims["user_id"].(float64)))

		return c.Next()
	}
}

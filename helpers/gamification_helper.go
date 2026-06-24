package helpers

import "math"

// CalculateXP mengembalikan jumlah XP berdasarkan tingkat kesulitan tugas.
// Nilai ini juga harus sinkron dengan konstanta xpOf() di main.go dan app.js.
func CalculateXP(difficulty string) int {
	switch difficulty {
	case "easy":
		return 10
	case "medium":
		return 20
	case "hard":
		return 30
	default:
		return 0
	}
}

// CalculateLevel menghitung level user dari total XP yang dimiliki.
// Formula: setiap 100 XP = 1 level, mulai dari level 1.
//   0–99 XP   → Level 1
//   100–199   → Level 2
//   200–299   → Level 3, dst.
func CalculateLevel(totalXP int) int {
	return int(math.Floor(float64(totalXP)/100)) + 1
}

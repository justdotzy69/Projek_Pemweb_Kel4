package helpers

import "math"

// CalculateXP menentukan jumlah XP berdasarkan tingkat kesulitan tugas
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

// CalculateLevel menentukan level pengguna saat ini berdasarkan total XP
// Asumsi: Setiap 100 XP akan naik 1 level. (0-99 = Lvl 1, 100-199 = Lvl 2, dst)
func CalculateLevel(totalXP int) int {
	return int(math.Floor(float64(totalXP)/100)) + 1
}
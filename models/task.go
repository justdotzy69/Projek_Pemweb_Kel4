package models

import "time"

// Task merepresentasikan tabel "tasks" di database.
type Task struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint       `gorm:"not null"                 json:"user_id"`

	// CategoryID pakai pointer (*uint) supaya bisa null (tugas tanpa kategori)
	CategoryID  *uint      `gorm:"default:null"             json:"category_id"`

	Title       string     `gorm:"not null"                 json:"title"`
	Description string     `gorm:"type:text"                json:"description"`
	Difficulty  string     `gorm:"type:enum('easy','medium','hard');not null" json:"difficulty"`
	Status      string     `gorm:"type:enum('pending','completed');default:'pending'" json:"status"`

	// Deadline juga pointer supaya bisa null (tugas tanpa deadline)
	Deadline    *time.Time `json:"deadline"`

	// Relasi: satu task punya satu kategori (opsional)
	Category    Category   `gorm:"foreignKey:CategoryID"    json:"category,omitempty"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

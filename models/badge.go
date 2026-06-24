package models

import "time"

// Badge merepresentasikan tabel "badges" di database.
// Badge diberikan ke user secara otomatis saat naik level di CompleteTask.
type Badge struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"unique;not null"          json:"name"`
	ImageURL      string    `gorm:"not null"                 json:"image_url"`

	// User akan mendapat badge ini ketika CurrentLevel mereka >= RequiredLevel
	RequiredLevel int       `gorm:"not null"                 json:"required_level"`

	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

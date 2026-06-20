package models

import "time"

type Badge struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"unique;not null" json:"name"`
	ImageURL      string    `gorm:"not null" json:"image_url"`
	RequiredLevel int       `gorm:"not null" json:"required_level"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
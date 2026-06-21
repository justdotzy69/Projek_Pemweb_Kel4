package models

import "time"

type Task struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint       `gorm:"not null" json:"user_id"`
	CategoryID  *uint      `gorm:"default:null" json:"category_id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Difficulty  string     `gorm:"type:enum('easy','medium','hard');not null" json:"difficulty"`
	Status      string     `gorm:"type:enum('pending','completed');default:'pending'" json:"status"`
	Deadline    *time.Time `json:"deadline"`
	Category    Category   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
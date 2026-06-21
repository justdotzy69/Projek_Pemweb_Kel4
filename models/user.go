package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email        string    `gorm:"unique;not null" json:"email"`
	Password     string    `gorm:"not null" json:"-"` // "-" menyembunyikan password dari response API
	Role         string    `gorm:"type:enum('user','admin');default:'user'" json:"role"`
	TotalXP      int       `gorm:"default:0" json:"total_xp"`
	CurrentLevel int       `gorm:"default:1" json:"current_level"`
	Tasks        []Task    `gorm:"constraint:OnDelete:CASCADE;" json:"tasks,omitempty"`
	Badges       []Badge   `gorm:"many2many:user_badges;constraint:OnDelete:CASCADE;" json:"badges,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
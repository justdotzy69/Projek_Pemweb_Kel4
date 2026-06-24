package models

import "time"

// User merepresentasikan tabel "users" di database.
// Field Password pakai json:"-" supaya tidak pernah ikut keluar di response API.
type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email        string    `gorm:"unique;not null"          json:"email"`
	Password     string    `gorm:"not null"                 json:"-"`
	Role         string    `gorm:"type:enum('user','admin');default:'user'" json:"role"`
	TotalXP      int       `gorm:"default:0"                json:"total_xp"`
	CurrentLevel int       `gorm:"default:1"                json:"current_level"`

	// Relasi: satu user punya banyak task. Kalau user dihapus, task ikut terhapus.
	Tasks  []Task  `gorm:"constraint:OnDelete:CASCADE;"             json:"tasks,omitempty"`

	// Relasi many-to-many dengan Badge lewat tabel pivot "user_badges"
	Badges []Badge `gorm:"many2many:user_badges;constraint:OnDelete:CASCADE;" json:"badges,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

package models

import "time"

// Category merepresentasikan tabel "categories" di database.
type Category struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"unique;not null"          json:"name"`

	// Kalau kategori dihapus, category_id di task di-set NULL (bukan ikut terhapus)
	Tasks []Task `gorm:"constraint:OnDelete:SET NULL;" json:"tasks,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

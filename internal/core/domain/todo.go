package domain

import "github.com/google/uuid"

type Todo struct {
	Model
	Title  string    `gorm:"not null" json:"title"`
	Done   bool      `gorm:"default:false" json:"done"`
	UserID uuid.UUID `gorm:"type:uuid" json:"user_id"`
}

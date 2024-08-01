package entities

import "time"

type Product struct {
	Id          uint      `gorm:"primaryKey"`
	Name        string
	CategoryId  uint      // Foreign key
	Category    Category  `gorm:"foreignKey:CategoryId"` // Relationship
	Stock       int64
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
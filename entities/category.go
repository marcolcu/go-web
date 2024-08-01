package entities

import "time"

type Category struct {
	Id        uint      `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
package model

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Id        string    `gorm:"primaryKey;type:uuid;"`
	Name      string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now()"`
}

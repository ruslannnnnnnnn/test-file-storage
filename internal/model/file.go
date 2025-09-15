package model

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

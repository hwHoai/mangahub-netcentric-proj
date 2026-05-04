package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	// BaseModel defines the basic structure and methods for all models.
	CreatedAt time.Time `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
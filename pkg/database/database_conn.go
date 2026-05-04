package database

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

type DatabaseConnection interface {
	InitDB(dbPath string) (*gorm.DB, error)
}
package database

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

type DatabaseConnectionInterface interface {
	InitDB(dbPath string) (*gorm.DB, error)
}
package main

import (
	"gorm.io/gorm"
	"mon-tool-be/models"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(models.Monitor{})
	db.AutoMigrate(models.Record{})
	db.AutoMigrate(models.Alert{})
}

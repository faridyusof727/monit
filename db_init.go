package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func InitDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB")
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

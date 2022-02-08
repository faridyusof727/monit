package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"mon-tool-be/handlers"
)

func initRouter(e *echo.Echo, db *gorm.DB) {
	// Monitor Routes
	e.POST("/monitors", handlers.Monitor{DB: db}.Store)
}

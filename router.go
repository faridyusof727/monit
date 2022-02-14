package main

import (
	"mon-tool-be/handlers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func initRouter(e *echo.Echo, db *gorm.DB) {
	// Monitor Routes
	e.POST("/monitors", handlers.Monitor{DB: db}.Store)
	e.GET("/monitors", handlers.Monitor{DB: db}.List)
	e.GET("/monitors/:id", handlers.Monitor{DB: db}.View)
	e.DELETE("/monitors/:id", handlers.Monitor{DB: db}.Delete)
	e.PATCH("/monitors/:id", handlers.Monitor{DB: db}.Edit)

	// Records
	e.GET("/monitors/:id/records", handlers.Record{DB: db}.View)
}

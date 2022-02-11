package main

import (
	"mon-tool-be/middlewares"
	"mon-tool-be/utils"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Load config
	err := godotenv.Load(".env")
	if err != nil {
		e.Logger.Fatal("Error loading .env file")
	}

	// Setting up timezone
	loc, err := time.LoadLocation(os.Getenv("APP_TZ"))
	if err != nil {
		e.Logger.Fatal("Error loading timezone")
	}
	time.Local = loc

	// Init DB
	db, err := InitDB()
	if err != nil {
		e.Logger.Fatal("Error connecting to DB")
	}

	// Loading telegram go routine
	go utils.InitTelegram(db)

	// DB Migration
	Migrate(db)

	// Init Cron
	InitCron(db)

	// Firebase instance
	firebaseAuth := utils.FirebaseAuth("firebase-admin-sdk.json")
	firebase := middlewares.Auth{
		AuthClient: firebaseAuth,
	}

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.Use(firebase.Check)

	// Router
	initRouter(e, db)

	// Server start
	err = e.Start(":1323")
	if err != nil {
		e.Logger.Fatal(err)
	}
}

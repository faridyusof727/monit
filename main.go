package main

import (
	echoSwagger "github.com/swaggo/echo-swagger"
	"mon-tool-be/middlewares"
	"mon-tool-be/utils"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "mon-tool-be/docs"
)

// @title Uptime Agent Swagger
// @version 1.0
// @description This is the API for Uptime Agent Swagger.

// @contact.name Farid Yusof
// @contact.url https://www.facebook.com/faridyusof727/
// @contact.email faridyusof727@gmail.com

// @license.name GNU GENERAL PUBLIC LICENSE
// @license.url https://github.com/faridyusof727/monit/blob/main/LICENSE

// @host localhost:1323
// @BasePath /
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
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Router
	g := e.Group("")
	g.Use(firebase.Check)
	initRouter(g, db)

	// Server start
	err = e.Start(":1323")
	if err != nil {
		e.Logger.Fatal(err)
	}
}

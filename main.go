package main

import (
	"log"
	"manajemen_tugas_master/app"
	"manajemen_tugas_master/helper"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

func main() {
	// Load environment variables
	helper.LoadEnv()

	// Connect to the database
	db, err := app.ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initialize Fiber
	fiberApp := fiber.New()

	// Konfigurasi CORS
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins:     "https://master.d3nck08c8eblbc.amplifyapp.com,http://127.0.0.1:5173",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization, GoogleAuthorization",
		ExposeHeaders:    "Content-Length,Set-Cookie,Authorization, GoogleAuthorization",
		AllowCredentials: true,
		MaxAge:           int((12 * time.Hour).Seconds()),
	}))

	// Middleware
	fiberApp.Use(logger.New())
	fiberApp.Use(recover.New())

	store = session.New()

	// Setup routes
	app.SetupRoutes(fiberApp, db, store)

	// Start Fiber
	log.Fatal(fiberApp.Listen(":" + os.Getenv("PORT")))
}

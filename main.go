package main

import (
	"log"
	"manajemen_tugas_master/app"
	"manajemen_tugas_master/helper"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

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

	// // Konfigurasi CORS
	// fiberApp.Use(cors.New(cors.Config{
	// 	AllowOrigins:     "http://127.0.0.1:5500",
	// 	AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
	// 	AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
	// 	ExposeHeaders:    "Content-Length",
	// 	AllowCredentials: true,
	// 	MaxAge:           int((12 * time.Hour).Seconds()),
	// }))

	// Middleware
	fiberApp.Use(logger.New())
	fiberApp.Use(recover.New())

	// Setup routes
	app.SetupRoutes(fiberApp, db)

	// Start Fiber
	log.Fatal(fiberApp.Listen(":" + os.Getenv("PORT")))
}

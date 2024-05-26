package main

import (
	"fmt"
	"manajemen_tugas_master/app"
	"manajemen_tugas_master/helper"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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

	// Middleware
	fiberApp.Use(logger.New())
	fiberApp.Use(recover.New())

	// Setup routes
	app.SetupRoutes(fiberApp, db)

	// cek .env
	port := os.Getenv("PORT")
	fmt.Println(port)

	// Start Fiber
	log.Fatal(fiberApp.Listen(":" + os.Getenv("PORT")))
}

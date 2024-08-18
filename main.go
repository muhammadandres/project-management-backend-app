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
	helper.LoadEnv()

	db, err := app.ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	fiberApp := fiber.New()

	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins:     "https://master.d3nck08c8eblbc.amplifyapp.com,http://127.0.0.1:5173,https://manajementugas.com,https://www.manajementugas.com",
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

	app.SetupRoutes(fiberApp, db, store)

	log.Fatal(fiberApp.Listen(":" + os.Getenv("PORT")))
}

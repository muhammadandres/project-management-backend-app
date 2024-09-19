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
	"github.com/gofiber/swagger"

	_ "manajemen_tugas_master/docs"
)

// @title           Project Management App
// @description     API documentation

// @contact.email  m.andres.novrizal@gmail.com
var store *session.Store

func main() {
	helper.LoadEnv()

	db, err := app.ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	fiberApp := fiber.New()
	fiberApp.Get("/swagger/*", swagger.HandlerDefault)

	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins:     "https://master.d3nck08c8eblbc.amplifyapp.com,http://127.0.0.1:5173,https://manajementugas.com,https://www.manajementugas.com",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization, GoogleAuthorization",
		ExposeHeaders:    "Content-Length,Set-Cookie,Authorization, GoogleAuthorization",
		AllowCredentials: true,
		MaxAge:           int((12 * time.Hour).Seconds()),
	}))

	fiberApp.Use(logger.New())
	fiberApp.Use(recover.New())

	// session
	store = session.New()

	app.SetupRoutes(fiberApp, db, store)

	// run code
	log.Fatal(fiberApp.Listen(":" + os.Getenv("PORT")))
}

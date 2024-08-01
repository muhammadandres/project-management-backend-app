package helper

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found or error loading it. Using environment variables directly.")
	} else {
		log.Println("Successfully loaded .env file")
	}
}

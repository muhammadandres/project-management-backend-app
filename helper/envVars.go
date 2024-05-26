package helper

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error in loading .env file: ", err)
	}
}

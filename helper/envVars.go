package helper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found or error loading it. Using environment variables directly.")
	} else {
		log.Println("Successfully loaded .env file")
	}

	requiredEnvVars := []string{"PORT", "DB_URL", "SECRET", "AWS_REGION", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "GOOGLE_CALENDAR_CREDENTIALS"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Required environment variable %s is not set", envVar)
		}
	}
}

package helper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found. Using environment variables directly.")
	} else {
		log.Println("Successfully loaded .env file")
	}

	requiredEnvVars := []string{"PORT", "DB_URL", "SECRET", "AWS_REGION", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "GOOGLE_CALENDAR_CREDENTIALS"}
	missingVars := []string{}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		log.Printf("Warning: The following required environment variables are not set: %v", missingVars)
		log.Println("Please ensure these variables are set in your environment or .env file.")
	} else {
		log.Println("All required environment variables are set.")
	}
}

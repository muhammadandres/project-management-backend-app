package helper

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func SetupGoogleAuth() (*oauth2.Config, error) {
	credentials := os.Getenv("OAUTH_CREDENTIALS")

	var credentialsJSON map[string]interface{}

	// Coba parse sebagai JSON object terlebih dahulu
	err := json.Unmarshal([]byte(credentials), &credentialsJSON)
	if err != nil {
		// Jika gagal, coba parse sebagai string JSON
		var credentialsString string
		err = json.Unmarshal([]byte(credentials), &credentialsString)
		if err != nil {
			return nil, fmt.Errorf("unable to parse OAUTH_CREDENTIALS: %v", err)
		}
		// Parse string JSON menjadi object
		err = json.Unmarshal([]byte(credentialsString), &credentialsJSON)
		if err != nil {
			return nil, fmt.Errorf("unable to parse OAUTH_CREDENTIALS content: %v", err)
		}
	}

	// Konversi kembali ke JSON bytes untuk google.ConfigFromJSON
	credentialsBytes, err := json.Marshal(credentialsJSON)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal credentials: %v", err)
	}

	config, err := google.ConfigFromJSON(credentialsBytes, "email", "profile")
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	return config, nil
}

func GenerateRandomState() string {
	// Fungsi ini tidak perlu diubah
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

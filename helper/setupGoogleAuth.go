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

	// Ekstrak informasi yang diperlukan dari credentialsJSON
	webConfig, ok := credentialsJSON["web"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid OAUTH_CREDENTIALS format: 'web' key not found or not an object")
	}

	clientID, ok := webConfig["client_id"].(string)
	if !ok {
		return nil, fmt.Errorf("client_id not found or not a string")
	}

	clientSecret, ok := webConfig["client_secret"].(string)
	if !ok {
		return nil, fmt.Errorf("client_secret not found or not a string")
	}

	redirectURIs, ok := webConfig["redirect_uris"].([]interface{})
	if !ok || len(redirectURIs) == 0 {
		return nil, fmt.Errorf("redirect_uris not found or empty")
	}

	redirectURL, ok := redirectURIs[0].(string)
	if !ok {
		return nil, fmt.Errorf("first redirect_uri is not a string")
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	return config, nil
}

func GenerateRandomState() string {
	// Fungsi ini tidak perlu diubah
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

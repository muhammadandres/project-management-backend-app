package helper

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var HardcodedCredentials = []byte(`{
    "web": {
        "client_id": "91088718933-p44u1h8q4n5s5hrbj6rkdbhk58tgfhce.apps.googleusercontent.com",
        "project_id": "main-crow-387504",
        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://oauth2.googleapis.com/token",
        "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
        "client_secret": "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",
        "redirect_uris": [
            "https://www.manajementugas.com/auth/callback"
        ]
    }
}`)

func SetupGoogleAuth() (*oauth2.Config, error) {
	// Use hardcoded credentials instead of reading from file
	var credentialsJSON map[string]interface{}
	err := json.Unmarshal(HardcodedCredentials, &credentialsJSON)
	if err != nil {
		return nil, fmt.Errorf("unable to parse hardcoded credentials: %v", err)
	}

	config, err := google.ConfigFromJSON(HardcodedCredentials,
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile")
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	config.RedirectURL = "https://www.manajementugas.com/auth/callback"
	return config, nil
}

func GenerateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

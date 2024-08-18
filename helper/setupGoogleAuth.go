package helper

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func SetupGoogleAuth() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
	return config
}

func GenerateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

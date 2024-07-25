package helper

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func SetupGoogleAuth() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "http://127.0.0.1:4040/auth/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
	return config
}
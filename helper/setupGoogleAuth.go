package helper

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func SetupGoogleAuth() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     "91088718933-p44u1h8q4n5s5hrbj6rkdbhk58tgfhce.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "https://www.manajementugas.com/auth/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	return config
}

func GenerateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

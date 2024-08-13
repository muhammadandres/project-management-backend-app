package helper

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
)

func VerifyTurnstileToken(token string) error {
	secretKey := os.Getenv("TURNSTILE_SECRET_KEY")
	endpoint := "https://challenges.cloudflare.com/turnstile/v0/siteverify"

	data := url.Values{}
	data.Set("secret", secretKey)
	data.Set("response", token)

	resp, err := http.PostForm(endpoint, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if success, ok := result["success"].(bool); !ok || !success {
		return errors.New("Turnstile verification failed")
	}

	return nil
}

package helper

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

type TurnstileResponse struct {
	Success     bool     `json:"success"`      // Menunjukkan apakah verifikasi berhasil
	ChallengeTs string   `json:"challenge_ts"` // Timestamp verifikasi tantangan diselesaikan
	Score       float64  `json:"score"`        // Skor kepercayaan yang menunjukkan seberapa yakin sistem bahwa interaksi dilakukan oleh manusia
	ErrorCodes  []string `json:"error-codes"`  // Array yang berisi code error jika ada masalah pada sistem
}

func VerifyTurnstileToken(token string) (*TurnstileResponse, error) {
	secretKey := os.Getenv("TURNSTILE_SECRET_KEY")
	resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify",
		url.Values{
			"secret":   {secretKey},
			"response": {token},
		})

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TurnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

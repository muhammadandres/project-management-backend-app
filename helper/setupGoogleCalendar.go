package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var hardcodedCredentials = []byte(`{
    "web": {
        "client_id": "91088718933-3mu4cb8n400hedbo9donc70ft7jjo90u.apps.googleusercontent.com",
        "project_id": "main-crow-387504",
        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://oauth2.googleapis.com/token",
        "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
        "client_secret": "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",
        "redirect_uris": [
            "https://www.manajementugas.com/auth/callback",
            "https://www.manajementugas.com"
        ]
    }
}`)

func createCalendarService(userEmail string) (*calendar.Service, error) {
	ctx := context.Background()

	// Use hardcoded credentials instead of reading from file
	b := hardcodedCredentials

	log.Printf("Contents of credentials: %s", string(b))

	// Parse the JSON configuration
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Printf("Error parsing credentials JSON: %v", err)
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	// Use the config to get a client
	client := getClient(config, userEmail)

	// Create the calendar service
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Error creating calendar service: %v", err)
		return nil, fmt.Errorf("unable to retrieve Calendar client: %v", err)
	}

	return srv, nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "rundll32"
		args = append(args, "url.dll,FileProtocolHandler")
	case "darwin":
		cmd = "open"
	default:
		return fmt.Errorf("unsupported platform")
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("Go to the following link in your browser: \n%v\n", authURL)

	err := openBrowser(authURL)
	if err != nil {
		log.Fatalf("Unable to open browser: %v", err)
	}

	codeChan := make(chan string)
	srv := &http.Server{Addr: "localhost:4040"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			fmt.Fprintf(w, "Authorization code received. You can close this window.")
			codeChan <- code
		} else {
			fmt.Fprintf(w, "No authorization code received.")
		}
	})
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Unable to start server: %v", err)
		}
	}()

	code := <-codeChan
	srv.Shutdown(context.TODO())

	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func getClient(config *oauth2.Config, userEmail string) *http.Client {
	tok, err := tokenFromFile()
	if err != nil {
		log.Printf("Error getting token from file: %v", err)
		tok, err = tokenFromEnv()
		if err != nil {
			log.Printf("Error getting token from env: %v", err)
			tok = getTokenFromWeb(config)
		}
	}

	// Check if token is expired
	if tok.Expiry.Before(time.Now()) {
		log.Println("Token has expired. Refreshing...")
		newTok, err := refreshToken(config, tok)
		if err != nil {
			log.Printf("Error refreshing token: %v", err)
			tok = getTokenFromWeb(config)
		} else {
			tok = newTok
		}
	}

	saveToken(tok)
	return config.Client(context.Background(), tok)
}

func tokenFromFile() (*oauth2.Token, error) {
	tokenFile := filepath.Join(".", "token.json")
	f, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func tokenFromEnv() (*oauth2.Token, error) {
	tokenJSON := os.Getenv("GOOGLE_OAUTH_TOKEN")
	if tokenJSON == "" {
		return nil, fmt.Errorf("GOOGLE_OAUTH_TOKEN not set")
	}
	tok := &oauth2.Token{}
	err := json.Unmarshal([]byte(tokenJSON), tok)
	return tok, err
}

func saveToken(token *oauth2.Token) error {
	// Save to file
	tokenFile := filepath.Join(".", "token.json")
	f, err := os.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		return err
	}

	// Save to environment variable
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return err
	}

	// Read .env file
	envContent, err := os.ReadFile(".env")
	if err != nil {
		return err
	}

	// Check if GOOGLE_OAUTH_TOKEN already exists
	lines := strings.Split(string(envContent), "\n")
	tokenFound := false
	for i, line := range lines {
		if strings.HasPrefix(line, "GOOGLE_OAUTH_TOKEN=") {
			lines[i] = fmt.Sprintf("GOOGLE_OAUTH_TOKEN=%s", string(tokenJSON))
			tokenFound = true
			break
		}
	}

	// If not found, add it at the end of the file
	if !tokenFound {
		lines = append(lines, fmt.Sprintf("GOOGLE_OAUTH_TOKEN=%s", string(tokenJSON)))
	}

	// Write back to .env file
	err = os.WriteFile(".env", []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return err
	}

	// Set environment variable
	return os.Setenv("GOOGLE_OAUTH_TOKEN", string(tokenJSON))
}

func refreshToken(config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error) {
	newToken, err := config.TokenSource(context.Background(), token).Token()
	if err != nil {
		return nil, err
	}
	if newToken.AccessToken != token.AccessToken {
		saveToken(newToken)
	}
	return newToken, nil
}

func CreateGoogleCalendarEvent(senderEmail, summary, description, startDateTime, endDateTime, timeZone string, attendees []string) (*calendar.Event, string, error) {
	// Use hardcoded credentials
	config, err := google.ConfigFromJSON(hardcodedCredentials, calendar.CalendarScope)
	if err != nil {
		return nil, "", fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	// Generate the authorization URL
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	srv, err := createCalendarService(senderEmail)
	if err != nil {
		return nil, authURL, err
	}

	event := &calendar.Event{
		Summary:     summary,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: startDateTime,
			TimeZone: timeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: endDateTime,
			TimeZone: timeZone,
		},
		Attendees: make([]*calendar.EventAttendee, 0, len(attendees)),
	}

	for _, email := range attendees {
		if email != senderEmail {
			event.Attendees = append(event.Attendees, &calendar.EventAttendee{Email: email})
		}
	}

	calendarId := "primary"
	event, err = srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		return nil, authURL, fmt.Errorf("unable to create event: %v", err)
	}

	return event, authURL, nil
}

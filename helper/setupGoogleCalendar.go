package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

	config, err := google.ConfigFromJSON(hardcodedCredentials, calendar.CalendarScope)
	if err != nil {
		log.Printf("Error parsing credentials JSON: %v", err)
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client := getClient(config, userEmail)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Error creating calendar service: %v", err)
		return nil, fmt.Errorf("unable to retrieve Calendar client: %v", err)
	}

	return srv, nil
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("Go to the following link in your browser: \n%v\n", authURL)
	fmt.Println("Enter the authorization code:")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

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
		tok = getTokenFromWeb(config)
		SaveToken(tok)
	}

	if tok.Expiry.Before(time.Now()) {
		log.Println("Token has expired. Refreshing...")
		newTok, err := refreshToken(config, tok)
		if err != nil {
			log.Printf("Error refreshing token: %v", err)
			tok = getTokenFromWeb(config)
		} else {
			tok = newTok
		}
		SaveToken(tok)
	}

	return config.Client(context.Background(), tok)
}

func tokenFromFile() (*oauth2.Token, error) {
	tokenFile := "/app/token.json"
	f, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func SaveToken(token *oauth2.Token) error {
	tokenFile := "/app/token.json"
	f, err := os.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

func refreshToken(config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error) {
	newToken, err := config.TokenSource(context.Background(), token).Token()
	if err != nil {
		return nil, err
	}
	if newToken.AccessToken != token.AccessToken {
		SaveToken(newToken)
	}
	return newToken, nil
}

func CreateGoogleCalendarEvent(senderEmail, summary, description, startDateTime, endDateTime, timeZone string, attendees []string) (*calendar.Event, string, error) {
	config, err := google.ConfigFromJSON(hardcodedCredentials, calendar.CalendarScope)
	if err != nil {
		return nil, "", fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

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

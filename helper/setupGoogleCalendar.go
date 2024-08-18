package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func createCalendarService(senderEmail string) (*calendar.Service, error) {
	credentials := os.Getenv("GOOGLE_CALENDAR_CREDENTIALS")

	var credentialsJSON map[string]interface{}

	// Try to parse as JSON object first
	err := json.Unmarshal([]byte(credentials), &credentialsJSON)
	if err != nil {
		// If failed, try to parse as JSON string
		var credentialsString string
		err = json.Unmarshal([]byte(credentials), &credentialsString)
		if err != nil {
			return nil, fmt.Errorf("unable to parse GOOGLE_CALENDAR_CREDENTIALS: %v", err)
		}
		// Parse JSON string to object
		err = json.Unmarshal([]byte(credentialsString), &credentialsJSON)
		if err != nil {
			return nil, fmt.Errorf("unable to parse GOOGLE_CALENDAR_CREDENTIALS content: %v", err)
		}
	}

	// Convert back to JSON bytes for ConfigFromJSON
	credentialsBytes, err := json.Marshal(credentialsJSON)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal credentials: %v", err)
	}

	config, err := google.ConfigFromJSON(credentialsBytes, calendar.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret to config: %v", err)
	}

	ctx := context.Background()
	client := getClient(config, "m.andres.novrizal@gmail.com")

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Calendar client: %v", err)
	}

	return srv, nil
}

func getClient(config *oauth2.Config, userEmail string) *http.Client {
	tokFile := fmt.Sprintf("token_%s.json", userEmail)
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		log.Printf("Error reading token file: %v", err)
		// Instead of getting a new token from the web, we'll use the existing token file
		tok, err = tokenFromFile("token_m.andres.novrizal@gmail.com.json")
		if err != nil {
			log.Fatalf("Unable to read token file: %v", err)
		}
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func CreateGoogleCalendarEvent(senderEmail, summary, description, startDateTime, endDateTime, timeZone string, attendees []string) (*calendar.Event, error) {
	srv, err := createCalendarService(senderEmail)
	if err != nil {
		return nil, err
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

	// Only add attendee if not the sender
	for _, email := range attendees {
		if email != senderEmail {
			event.Attendees = append(event.Attendees, &calendar.EventAttendee{Email: email})
		}
	}

	calendarId := "primary"
	event, err = srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to create event: %v", err)
	}

	return event, nil
}

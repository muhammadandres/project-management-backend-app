package helper

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendEmail(to []string, subject, body string) error {
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	smtpHost := os.Getenv("SMTPHOST")
	smtpPort := os.Getenv("SMTPPORT")

	from := "m.andres.novrizal@gmail.com"

	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to[0], from, subject, body))

	auth := smtp.PlainAuth("", username, password, smtpHost)
	smtpAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	log.Printf("Attempting to send email via %s", smtpAddr)
	err := smtp.SendMail(smtpAddr, auth, from, to, message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Println("Email sent successfully")
	return nil
}

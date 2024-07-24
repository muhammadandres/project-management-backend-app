package helper

import (
	"fmt"
	"net/smtp"

	"github.com/gofiber/fiber/v2"
)

func SendEmail(c *fiber.Ctx, to []string, subject, body string) error {
	username := "api"
	password := "your_mailtrap_password"
	smtpHost := "live.smtp.mailtrap.io"
	smtpPort := "587"

	auth := smtp.PlainAuth("", username, password, smtpHost)

	from := "anyname@freelance.mailtrap.link"

	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to[0], from, subject, body))

	smtpUrl := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err := smtp.SendMail(smtpUrl, auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}

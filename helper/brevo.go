package helper

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendEmail(to []string, subject, htmlBody string) error {
	username := os.Getenv("BREVO_USERNAME")
	password := os.Getenv("BREVO_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMPTP_PORT")

	from := "manajementugasapp@gmail.com"

	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", to[0], from, subject, htmlBody))

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

func GetEmailTemplate(title, taskName, status, description string) string {
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>%s</title>
        <style>
            body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
            .container { max-width: 600px; margin: 0 auto; padding: 20px; }
            .header { background-color: #f4f4f4; padding: 10px; text-align: center; }
            .content { padding: 20px; background-color: #ffffff; }
            .footer { text-align: center; padding: 10px; font-size: 0.8em; color: #777; }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>%s</h1>
            </div>
            <div class="content">
                <p>Hello,</p>
                <p>We're writing to inform you that the task "%s" has been updated.</p>
                <p><strong>Status:</strong> %s</p>
                <p>%s</p>
                <p>If you have any questions or need further information, please don't hesitate to contact m.andres.novrizal@gmail.com</p>
            </div>
            <div class="footer">
                <p>This is an automated message. Please do not reply directly to this email.</p>
            </div>
        </div>
    </body>
    </html>
    `, title, title, taskName, status, description)
}

func GetCalendarInviteTemplate(summary, description string) string {
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Calendar Invite Notification</title>
        <style>
            body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
            .container { max-width: 600px; margin: 0 auto; padding: 20px; }
            .header { background-color: #4285f4; color: white; padding: 10px; text-align: center; }
            .content { padding: 20px; background-color: #ffffff; }
            .button { display: inline-block; padding: 10px 20px; background-color: #4285f4; color: white; text-decoration: none; border-radius: 5px; }
            .footer { text-align: center; padding: 10px; font-size: 0.8em; color: #777; }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>Calendar Invite Notification</h1>
            </div>
            <div class="content">
                <h2>%s</h2>
                <p>You have been invited to an event. Here are the details:</p>
                <p><strong>Description:</strong> %s</p>
                <p><strong>Important Notes:</strong></p>
                <ul>
                    <li>The official calendar invitation will be sent separately by the Google Calendar system shortly.</li>
                    <li>If you haven't received the invitation yet, please wait a few moments and check your inbox periodically.</li>
                    <li>Once you receive the invitation, don't forget to click the "Add to calendar" button on the invite.</li>
                </ul>
                <p>Please check your calendar in the next few minutes to see the date and time of this event.</p>
            </div>
            <div class="footer">
                <p>This is an automated message. Please do not reply directly to this email.</p>
            </div>
        </div>
    </body>
    </html>
    `, summary, description)
}

func ForgotPasswordTemplate(resetCode string) string {
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Password Reset Code</title>
        <style>
            body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; background-color: #f4f4f4; }
            .container { max-width: 600px; margin: 20px auto; padding: 20px; background-color: #ffffff; border-radius: 5px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
            .header { text-align: center; margin-bottom: 20px; }
            .content { padding: 20px; background-color: #ffffff; }
            .code { font-size: 32px; font-weight: bold; text-align: center; letter-spacing: 5px; margin: 20px 0; color: #007bff; }
            .footer { text-align: center; margin-top: 20px; font-size: 0.8em; color: #777; }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h2>Password Reset Code</h2>
            </div>
            <div class="content">
                <p>Your password reset code is:</p>
                <div class="code">%s</div>
                <p>This code will expire in 15 minutes.</p>
                <p>If you didn't request a password reset, please ignore this email or contact support if you have concerns.</p>
            </div>
            <div class="footer">
                <p>This is an automated message. Please do not reply directly to this email.</p>
            </div>
        </div>
    </body>
    </html>
    `, resetCode)
}

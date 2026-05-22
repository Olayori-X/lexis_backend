package functions

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func SendSimpleMessage(toEmail string, subject string, body string) (string, error) {
	apiKey := os.Getenv("MAILGUN_API_KEY")
	domain := os.Getenv("MAILGUN_DOMAIN")

	if apiKey == "" || domain == "" {
		return "", errors.New("mailgun credentials not set")
	}

	mg := mailgun.NewMailgun(domain, apiKey)

	message := mg.NewMessage(
		"Mailgun Sandbox <postmaster@"+domain+">",
		subject,
		body,
		toEmail,
	)

	log.Println("Sending email via Mailgun...")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, id, err := mg.Send(ctx, message)
	return id, err
}

func SendMail(to []string, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD") // use app password here

	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return fmt.Errorf("failed to send mail: %v", err)
	}

	return nil
}

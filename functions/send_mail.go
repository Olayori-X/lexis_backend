package functions

import "os"

func SendEmail(to, subject, body string) error {
	if os.Getenv("APP_ENV") == "production" {
		_, err := SendSimpleMessage(to, subject, body)
		return err
	}
	return SendMail([]string{to}, subject, body)
}

package pkg

import (
	"fmt"
	"os"
	"strconv"
	"subscriptionplus/server/infra/logger"

	"gopkg.in/gomail.v2"
)

// SendEmail отправка email письма на почту
func SendEmail(to, title, content string) error {
	logger := logger.NewLogger()

	go_env := os.Getenv("GO_ENV") == "DEV"

	mail := gomail.NewMessage()
	mail.SetHeader("From", os.Getenv("SMTP_EMAIL"))
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", title)
	mail.SetBody("text/html", content)

	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	d := gomail.NewDialer(os.Getenv("SMTP_ADDR"), port, os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"))

	if go_env {
		d.SSL = false
	} else {
		d.SSL = true
	}

	if err := d.DialAndSend(mail); err != nil {
		return fmt.Errorf("%s", "failed to send email")
	}

	logger.Info("email sent: %s", to)
	return nil
}

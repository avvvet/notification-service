// internal/email/email.go
package email

import (
	"fmt"
	"log"
	"net/smtp"
)

// Email struct represents an email.
type Email struct {
	Subject string
	Body    string
	To      string
}

// SendEmail sends the email using the provided SMTP configuration.
func SendEmail(email *Email, smtpConfig SMTPConfig) error {
	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)

	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", email.Subject, email.Body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", smtpConfig.Host, smtpConfig.Port),
		auth,
		smtpConfig.From,
		[]string{email.To},
		[]byte(msg),
	)

	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}

// SMTPConfig represents SMTP server configuration.
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

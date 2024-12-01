package smtp

import (
	"fmt"
	"net/smtp"

	"Questify/config"
)

// SMTPClient handles sending emails via an SMTP server
type SMTPClient struct {
	host       string
	port       int
	username   string
	password   string
	senderEmail string
}

// NewSMTPClient creates a new SMTPClient instance
func NewSMTPClient(cfg config.SMTP) *SMTPClient {
	return &SMTPClient{
		host:       cfg.Host,
		port:       cfg.Port,
		username:   cfg.Username,
		password:   cfg.Password,
		senderEmail: cfg.SenderEmail,
	}
}

// SendEmail sends an email using the SMTP server
func (s *SMTPClient) SendEmail(recipient, subject, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\n\r\n%s",
		recipient, subject, body,
	))

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	if err := smtp.SendMail(addr, auth, s.senderEmail, []string{recipient}, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

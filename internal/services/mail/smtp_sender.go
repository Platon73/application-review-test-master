package mail

import (
	"fmt"
	"net/smtp"
)

type SMTPSender struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func (s *SMTPSender) SendConfirmation(to, content string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Booking Confirmation\r\n\r\n%s",
		s.From, to, content,
	)
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", s.Host, s.Port),
		auth,
		s.From,
		[]string{to},
		[]byte(msg),
	)
}

package service

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/internal/config"
)

type emailService struct {
	cnf *config.Config
}

func NewEmailService(cnf *config.Config) domain.EmailService {
	return &emailService{cnf: cnf}
}

// Send implements domain.EmailService.
// func (s emailService) Send(to string, subject string, body string) error {
// 	auth := smtp.PlainAuth("", s.cnf.Email.User, s.cnf.Email.Password, s.cnf.Email.Host)
// 	msg := []byte("From: Green Basket <" + s.cnf.Email.User + ">\n" +
// 		"To: " + to + "\n" +
// 		"Subject: " + subject + "\n" +
// 		body)

// 	return smtp.SendMail(s.cnf.Email.Host+":"+s.cnf.Email.Port, auth, s.cnf.Email.User, []string{to}, msg)
// }

func (s emailService) SendMail(to, subject, body string) error {
	auth := smtp.PlainAuth(
		"",
		s.cnf.Email.Name,
		s.cnf.Email.Password,
		s.cnf.Email.Host,
	)

	msg := "Subject: " + subject + "\n" + body
	fullHost := fmt.Sprintf("%s:%d", s.cnf.Email.Host, 587)
	err := smtp.SendMail(
		fullHost,
		auth,
		s.cnf.Email.Name,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	return nil
}

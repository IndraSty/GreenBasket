package domain

type EmailService interface {
	SendMail(to, subject, body string) error
}

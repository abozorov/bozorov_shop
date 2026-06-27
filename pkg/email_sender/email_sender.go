package emailsender

import (
	"net/smtp"
)

type EmailSender struct {
	from     string
	password string
	smtpHost string
	smtpPort string
	auth     smtp.Auth
}

func NewEmailSender(from, password, smtpHost, smtpPort string) *EmailSender {

	return &EmailSender{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		auth:     smtp.PlainAuth("", from, password, smtpHost),
	}
}

func (e *EmailSender) SendMail(to, message string) error {
	subject := "Subject: Mail Autentification!\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	body := "Код для верификации: " + message

	sendMess := []byte(subject + mime + body)

	return smtp.SendMail(e.smtpHost+":"+e.smtpPort, e.auth, e.from, []string{to}, sendMess)
}

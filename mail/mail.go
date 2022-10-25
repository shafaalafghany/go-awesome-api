package mailer

import (
	"fmt"
	"log"
	"net/smtp"
)

const sender = "eLibrary %3cno-reply@elibrary.com%3e"

type Config struct {
	AppUrl       string
	MailHost     string
	MailPort     int
	MailUsername string
	MailPassword string
	MailSender   string
}

type Mailer struct {
	config *Config
}

type EmailSender interface {
	SendActivationLink(id int, recipient, content string)
}

func NewMail(cfg *Config) EmailSender {
	cfg.MailSender = sender
	m := &Mailer{
		config: cfg,
	}
	return m
}

func (m *Mailer) SendActivationLink(id int, recipient, content string) {
	from := "elibrary"
	to := []string{recipient}
	link := fmt.Sprintf("%s/auth/verify/id=%d&token=%s", m.config.AppUrl, id, content)
	msg := []byte(fmt.Sprintf("From : %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", recipient) +
		"Subject: Email Verification elibrary\r\n\r\n" + activationTemplate(link))
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", m.config.MailHost, m.config.MailPort),
		smtp.PlainAuth("", m.config.MailUsername, m.config.MailPassword, m.config.MailHost),
		m.config.MailSender,
		to,
		msg,
	)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func activationTemplate(link string) string {
	greet := "Hi There,\n\n"
	instruction := "Please activate you email by clicking the link below\n"
	regards := "\n\nCheers\nelibrary team"
	return fmt.Sprintf(greet+instruction+"%s"+regards, link)
}

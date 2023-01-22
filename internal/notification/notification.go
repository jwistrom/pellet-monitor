package notification

import (
	"log"
	"net/smtp"
	"strconv"
)

type Service interface {
	SendNotification(text string)
}

type EmailService struct {
	from       string
	address    string
	auth       smtp.Auth
	recipients []string
}

func NewEmailService(host string, port int, from string, password string, recipients []string) *EmailService {
	address := host + ":" + strconv.Itoa(port)
	auth := smtp.PlainAuth("", from, password, host)

	return &EmailService{
		from:       from,
		address:    address,
		auth:       auth,
		recipients: recipients,
	}
}

func (es *EmailService) SendNotification(text string) {
	subject := "Subject: Pelletsvarning!!\n"
	message := []byte(subject + text)

	err := smtp.SendMail(es.address, es.auth, es.from, es.recipients, message)
	if err != nil {
		log.Println("Failed to send warning", err)
	}
}

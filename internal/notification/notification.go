package notification

import (
	"log"
	"net/mail"
	"net/smtp"
	"strconv"
)

type Service interface {
	SendNotification(text string)
	AddRecipient(recipient string)
	DeleteRecipient(recipient string)
	GetRecipients() []string
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

func (es *EmailService) GetRecipients() []string {
	return es.recipients
}

func (es *EmailService) AddRecipient(recipient string) {
	_, err := mail.ParseAddress(recipient)
	if err != nil {
		log.Println("Not a valid email "+recipient+". ", err)
		return
	}
	es.recipients = append(es.recipients, recipient)
}

func (es *EmailService) DeleteRecipient(recipient string) {
	idxToDelete := es.indexOfRecipient(recipient)
	if idxToDelete != -1 {
		es.recipients = removeIndex(es.recipients, idxToDelete)
	} else {
		log.Println("No such recipient to delete " + recipient)
	}
}

func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func (es *EmailService) indexOfRecipient(recipient string) int {
	idx := -1
	for i, rec := range es.recipients {
		if rec == recipient {
			idx = i
			break
		}
	}
	return idx
}

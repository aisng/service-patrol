package main

import (
	"net/smtp"
	"os"
)

func SendMail(mailingList []string, message string) error {
	sender := os.Getenv("MAILUSERNAME")
	token := os.Getenv("MAILTOKEN")

	auth := smtp.PlainAuth(
		"",
		sender,
		token,
		"smtp.gmail.com",
	)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		sender,
		mailingList,
		[]byte(message),
	)
	if err != nil {
		return err
	}
	return nil
}

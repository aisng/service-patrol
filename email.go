package main

import (
	"net/smtp"
	"os"
)

func SendMail(mailingList []string, message string) error {
	sender := os.Getenv("SPMAILUSERNAME")
	token := os.Getenv("SPMAILTOKEN")

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

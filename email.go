package main

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendMail(mailingList []string, message string) {
	auth := smtp.PlainAuth(
		"",
		"***REMOVED***",
		os.Getenv("MAILTOKEN"),
		"smtp.gmail.com",
	)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"***REMOVED***",
		mailingList,
		[]byte(message),
	)
	if err != nil {
		fmt.Println(err)
	}
}

package main

import (
	"fmt"
	"net/smtp"
)

func sendMail(mailingList []string) {
	auth := smtp.PlainAuth(
		"",
		"***REMOVED***",
		"***REMOVED***",
		"smtp.gmail.com",
	)

	msg := "Subject: test subj2\ntest body2"

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"***REMOVED***",
		mailingList,
		[]byte(msg),
	)
	if err != nil {
		fmt.Println(err)
	}
}

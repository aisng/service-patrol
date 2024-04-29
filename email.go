package main

import (
	"fmt"
	"net/smtp"
	"os"
)

func getMessage(downServices []string, recoveredServices []string, nextCheckIn uint) string {
	var subject string
	var body string

	var downList string
	var recoveredList string

	nextCheckString := fmt.Sprintf("Next check will be made after %d hours.", nextCheckIn)

	areServicesDown := len(downServices) > 0
	areServicesRecovered := len(recoveredServices) > 0

	if areServicesDown {
		for _, service := range downServices {
			downList += " - " + service + "\n"
		}
	}

	if areServicesRecovered {
		for _, service := range recoveredServices {
			recoveredList += " - " + service + "\n"
		}
	}

	if areServicesDown && areServicesRecovered {
		subject = "Connection to FMC services recovered"
		body = fmt.Sprintf("Hello,\n\nconnection to the pages/IP's below are recovered:\n%s\nThe following pages are still down:\n%s\n%s", recoveredList, downList, nextCheckString)

	} else if areServicesDown && !areServicesRecovered {
		subject = "Connection to FMC services lost"
		body = fmt.Sprintf("Hello,\n\nconnection to the pages/IP's below are down:\n%s\n%s", downList, nextCheckString)

	} else if areServicesRecovered && !areServicesDown {
		subject = "Connection to FMC services recovered"
		body = fmt.Sprintf("Hello,\n\nconnection to the pages/IP's below are recovered:\n%s\n%s", recoveredList, nextCheckString)
	}

	return fmt.Sprintf("Subject: %s\n%s", subject, body)
}

func sendMail(mailingList []string, message string) {
	auth := smtp.PlainAuth(
		"",
		"***REMOVED***",
		// "***REMOVED***",
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

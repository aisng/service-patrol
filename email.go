package main

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

const messageTemplate string = `Subject: {{.Subject}}
Hello,
	
connection to the pages/IPs below was {{.GeneralStatus}}:
{{.GeneralList}}{{if .AdditionalList}}
The following pages are still down:
{{.AdditionalList}}{{end}}
Next check will be made after {{.Frequency}} hours.`

type Message struct {
	Subject        string
	GeneralStatus  string
	GeneralList    string
	AdditionalList string
	Frequency      uint
}

func NewMessage(downServices, recoveredServices []string, nextCheckIn uint) *Message {
	var subject string
	var generalStatus string
	var generalList string
	var additionalList string

	areServicesDown := len(downServices) > 0
	areServicesRecovered := len(recoveredServices) > 0

	if areServicesDown && areServicesRecovered {
		subject = "Connection to some FMC services recovered"
		generalStatus = "recovered"
		generalList = generateServicesList(recoveredServices)
		additionalList = generateServicesList(downServices)

	} else if areServicesDown && !areServicesRecovered {
		subject = "Connection to FMC services lost"
		generalStatus = "lost"
		generalList = generateServicesList(downServices)

	} else if areServicesRecovered && !areServicesDown {
		subject = "Connection to FMC services recovered"
		generalStatus = "recovered"
		generalList = generateServicesList(recoveredServices)
	}

	return &Message{
		Subject:        subject,
		GeneralStatus:  generalStatus,
		GeneralList:    generalList,
		AdditionalList: additionalList,
		Frequency:      nextCheckIn,
	}
}

func ParseTemplate(message Message) string {
	var output bytes.Buffer

	msgTmpl := template.Must(template.New("").Parse(messageTemplate))
	err := msgTmpl.Execute(&output, message)

	if err != nil {
		fmt.Println(err)
	}

	return output.String()
}

func generateServicesList(services []string) string {
	var list string
	for _, service := range services {
		list += " - " + service + "\n"
	}
	return list
}

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

package main

import (
	"bytes"
	"errors"
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

	if !(areServicesDown || areServicesRecovered) {
		return nil
	}

	if areServicesDown && areServicesRecovered {
		subject = "Connection to some FMC services recovered"
		generalStatus = "recovered"
		generalList = formatServicesListStr(recoveredServices)
		additionalList = formatServicesListStr(downServices)

	} else if areServicesDown && !areServicesRecovered {
		subject = "Connection to FMC services lost"
		generalStatus = "lost"
		generalList = formatServicesListStr(downServices)

	} else if areServicesRecovered && !areServicesDown {
		subject = "Connection to FMC services recovered"
		generalStatus = "recovered"
		generalList = formatServicesListStr(recoveredServices)
	}

	return &Message{
		Subject:        subject,
		GeneralStatus:  generalStatus,
		GeneralList:    generalList,
		AdditionalList: additionalList,
		Frequency:      nextCheckIn,
	}
}

func ParseTemplate(message *Message, templStr string) (string, error) {
	var output bytes.Buffer

	// if reflect.DeepEqual(message, &Message{}) {
	// 	return "", errors.New("struct is empty: nothing to parse")
	// }

	if message == nil {
		return "", errors.New("struct is empty: nothing to parse")
	}

	msgTmpl := template.Must(template.New("").Parse(templStr))
	err := msgTmpl.Execute(&output, message)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func formatServicesListStr(services []string) string {
	var list string
	for _, service := range services {
		list += " - " + service + "\n"
	}
	return list
}

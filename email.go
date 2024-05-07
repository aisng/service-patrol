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
	
connection to the pages/IPs below is {{.GeneralStatus}}:
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

func (m *Message) GenerateMessage(downServices, recoveredServices []string, nextCheckIn uint) {

	areServicesDown := len(downServices) > 0
	areServicesRecovered := len(recoveredServices) > 0

	if areServicesDown && areServicesRecovered {
		m.Subject = "Connection to some FMC services recovered"
		m.GeneralStatus = "recovered"
		m.GeneralList = generateServicesList(recoveredServices)
		m.AdditionalList = generateServicesList(downServices)

	} else if areServicesDown && !areServicesRecovered {
		m.Subject = "Connection to FMC services lost"
		m.GeneralStatus = "lost"
		m.GeneralList = generateServicesList(downServices)

	} else if areServicesRecovered && !areServicesDown {
		m.Subject = "Connection to FMC services recovered"
		m.GeneralStatus = "recovered"
		m.GeneralList = generateServicesList(recoveredServices)
	}
	m.Frequency = nextCheckIn
}

func (m *Message) ParseTemplate() string {
	var output bytes.Buffer

	msgTmpl := template.Must(template.New("messageTemplate").Parse(messageTemplate))
	err := msgTmpl.Execute(&output, m)

	if err != nil {
		fmt.Println(err)
	}

	return output.String()
}

// m = []byte(m)

// func getMessage(downServices []string, recoveredServices []string, nextCheckIn uint) string {
// 	var subject string
// 	var body string

// 	var downList string
// 	var recoveredList string

// 	nextCheckString := fmt.Sprintf("Next check will be made after %d hours.", nextCheckIn)

// 	areServicesDown := len(downServices) > 0
// 	areServicesRecovered := len(recoveredServices) > 0

// 	if areServicesDown {
// 		downList = generateServicesList(downServices)
// 	}
// 	ata := struct {
// 		GeneralStatus string
// 		GeneralList   string
// 		ExtraList     []string
// 		Frequency     uint
// 	}{
// 		GeneralStatus: "good",
// 		GeneralList:   "a, b, c, d",
// 		ExtraList:     []string{"f", "g", "h"},
// 		Frequency:     4,
// 	}

// 	if areServicesRecovered {
// 		recoveredList = generateServicesList(recoveredServices)
// 	}

// if areServicesDown && areServicesRecovered {
// 	subject = "Connection to some FMC services recovered"
// 	body = fmt.Sprintf("Hello,\n\nconnection to the pages/IPs below is recovered:\n%s\n"+
// 		"The following pages are still down:\n%s\n%s", recoveredList, downList, nextCheckString)

// } else if areServicesDown && !areServicesRecovered {
// 	subject = "Connection to FMC services lost"
// 	body = fmt.Sprintf("Hello,\n\nconnection to the pages/IPs below is down:\n%s\n%s", downList, nextCheckString)

// } else if areServicesRecovered && !areServicesDown {
// 	subject = "Connection to FMC services recovered"
// 	body = fmt.Sprintf("Hello,\n\nconnection to the pages/IPs below is recovered:\n%s\n%s", recoveredList, nextCheckString)
// }

// 	return fmt.Sprintf("Subject: %s\n%s", subject, body)
// }

// func generateSubject(isDownFound bool, isRecoveredFound bool) string {
// 	if isDownFound && isRecoveredFound {
// 		return "Connection to some FMC services recovered"
// 	}
// 	if isDownFound {
// 		return "Connection to FMC services lost"
// 	}
// 	if isRecoveredFound {
// 		return "Connection to FMC services recovered"
// 	}
// }

func generateServicesList(services []string) string {
	var list string
	for _, service := range services {
		list += " - " + service + "\n"
	}
	return list
}

func sendMail(mailingList []string, message string) {
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
	fmt.Println(mailingList)
}

package main

import (
	"log"
)

func main() {
	var config Config
	var status Status

	if err := config.Read(configFilename); err != nil {
		log.Fatalf("failed to read config file: %q: %v", configFilename, err)
	}

	if err := status.Read(statusFilename); err != nil {
		log.Printf("failed to read status file: %q: %v", statusFilename, err)
	}

	sp := NewServicePatrol(&config, &status)

	down, recovered, err := sp.Start()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if sp.IsDownLimitExceeded() || sp.IsRecoveredFound() {
		msg := NewMessage(down, recovered, config.Frequency)

		msgStr, err := ParseTemplate(msg, messageTemplate)
		if err != nil {
			log.Fatalf("error parsing template: %v", err)
		}

		err = SendMail(config.MailingList, msgStr)
		if err != nil {
			log.Fatalf("error sending mail:  %v", err)
		}

		if sp.IsDownLimitExceeded() {
			log.Printf("%d services are down (limit <= %d): email sent\n", status.DownCount, config.DownLimit)
		}

		if sp.IsRecoveredFound() {
			log.Printf("%d services recovered: email sent\n", len(recovered))
		}

	} else {
		log.Println("all services are running")
	}
}

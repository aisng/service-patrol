package main

import (
	"log"
	"time"
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

	timeChecked := time.Now().Local().Format("15:04:05")

	if sp.IsDownLimitExceeded() || sp.IsRecoveredFound() {
		msg := NewMessage(down, recovered, config.Frequency, timeChecked)

		msgStr, err := ParseTemplate(msg, messageTemplate)
		if err != nil {
			log.Fatalf("error parsing template: %v", err)
		}

		// fmt.Println(msgStr)
		err = SendMail(config.MailingList, msgStr)
		if err != nil {
			log.Fatalf("error sending mail:  %v", err)
		}

		if sp.IsDownLimitExceeded() {
			log.Printf("%d service(s) down. limit (%d) reached: email sent\n", status.DownCount, config.DownLimit)
		}

		if sp.IsRecoveredFound() {
			log.Printf("%d service(s) recovered: email sent\n", len(recovered))
		}
	}

	if status.DownCount > 0 && !sp.IsDownLimitExceeded() {
		log.Printf("%d service(s) down. limit (%d) not reached: email not sent\n", status.DownCount, config.DownLimit)
	}

	if status.DownCount == 0 {
		log.Println("all services are running")
	}
}

package main

import (
	"fmt"
)

func main() {
	var config Config
	var status Status

	if err := config.Read(configFilename); err != nil {
		panic(err)
	}

	if err := status.Read(statusFilename); err != nil {
		fmt.Println(err)
		return
	}

	sp := NewServicePatrol(&config, &status)

	down, recovered, err := sp.Start()
	if err != nil {
		fmt.Println(err)
	}

	if sp.IsDownLimitExceeded() || sp.IsRecoveredFound() {
		// TODO: figure out "chained" ptrs/deref
		msg := NewMessage(down, recovered, config.Frequency)
		msgStr, err := ParseTemplate(msg, messageTemplate)

		if err != nil {
			fmt.Println(err)
		}

		err = SendMail(config.MailingList, msgStr)
		if err != nil {
			panic(fmt.Errorf("error sending mail:  %v", err))
		}

		if sp.IsDownLimitExceeded() {
			fmt.Printf("%d services are down (limit <= %d): email sent\n", status.DownCount, config.DownLimit)
		}

		if sp.IsRecoveredFound() {
			fmt.Printf("%d services recovered: email sent\n", len(recovered))
		}
		// fmt.Println(msgStr)
	} else {
		fmt.Println("all services are running")
	}
}

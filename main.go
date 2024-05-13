package main

import (
	"fmt"
)

func main() {
	var config Config
	var serviceStatus ServiceStatus

	if err := config.Read(); err != nil {
		panic(err)
	}

	if err := serviceStatus.Read(); err != nil {
		fmt.Println(err)
		return
	}

	sp := NewServicePatrol(&config, &serviceStatus)

	down, recovered := sp.Start()

	if down != nil || recovered != nil {
		// TODO: figure out "chained" ptrs/deref
		msg := NewMessage(*down, *recovered, config.Frequency)
		msgStr, err := ParseTemplate(msg, messageTemplate)
		if err != nil {
			fmt.Println(err)
		}
		err = SendMail(config.MailingList, msgStr)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("down_count (%d) >= down_limit (%d). Email sent.\n", serviceStatus.DownCount, config.DownLimit)

		fmt.Println(msgStr)
	}
}

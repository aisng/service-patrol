package main

import (
	"fmt"
	"net/http"
	"time"
)

const (
	configFilename        string = "config.yaml"
	serviceStatusFilename string = "service-status.yaml"
)

func main() {
	var config Config
	var client http.Client
	var serviceStatus ServiceStatus
	var recoveredServices []string
	var affectedServices []string

	if err := config.Read("c" + configFilename); err != nil {
		panic(err)
		// return
	}

	if err := serviceStatus.Read(serviceStatusFilename); err != nil {
		fmt.Println(err)
		return
	}

	client = http.Client{
		Timeout: time.Second * time.Duration(config.Timeout),
	}

	for _, serviceUrl := range config.Services {
		isRunning, err := isServiceRunning(&client, serviceUrl)

		if err != nil {
			fmt.Println(err)
		}

		if !isRunning {
			affectedServices = append(affectedServices, serviceUrl)
			serviceStatus.incrementDownCount()
		}

		if isRunning && serviceStatus.isAffected(serviceUrl) {
			recoveredServices = append(recoveredServices, serviceUrl)
			serviceStatus.decrementDownCount()
		}
	}

	serviceStatus.AffectedServices = affectedServices

	if err := serviceStatus.Write(serviceStatusFilename); err != nil {
		fmt.Println(err)
	}

	isDownLimitExceeded := serviceStatus.DownCount >= config.DownLimit
	areServicesRecovered := len(recoveredServices) > 0

	if !isDownLimitExceeded {
		fmt.Printf("down_count (%d) <= down_limit (%d). Email will not be sent.\n", serviceStatus.DownCount, config.DownLimit)
	}

	if isDownLimitExceeded || areServicesRecovered {
		// TODO: figure out "chained" ptrs/deref
		msg := NewMessage(serviceStatus.AffectedServices, recoveredServices, config.Frequency)
		msgStr := ParseTemplate(msg)
		// SendMail(config.MailingList, msgStr)
		fmt.Printf("down_count (%d) >= down_limit (%d). Email sent.\n", serviceStatus.DownCount, config.DownLimit)

		fmt.Println(msgStr)
	}
}

func isServiceRunning(client *http.Client, url string) (bool, error) {
	resp, err := client.Head(url)
	if err != nil {
		return false, err
	}
	resp.Body.Close()
	return true, nil
}

package main

import (
	"fmt"
	"net/http"
	"time"
)

const (
	defaultTimeout        uint   = 3
	defaultLimit          uint   = 2
	defaultFrequency      uint   = 2
	configFilename        string = "config.yaml"
	serviceStatusFilename string = "service-status.yaml"
)

func main() {
	var config Config
	var client http.Client
	var serviceStatus ServiceStatus
	var recoveredServices []string
	var affectedServices []string

	if err := config.Read(configFilename); err != nil {
		fmt.Println(err)
		return
	}

	// if err := ReadYaml("c"+configFilename, config); err != nil {
	// 	fmt.Println(err)
	// }
	fmt.Println(config)

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

	if isDownLimitExceeded || areServicesRecovered {
		msg := NewMessage(serviceStatus.AffectedServices, recoveredServices, config.Frequency)
		msgStr := ParseTemplate(*msg)
		// sendMail(config.MailingList, msgStr)
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

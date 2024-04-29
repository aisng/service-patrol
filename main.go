package main

import (
	"fmt"
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
	var serviceStatus ServiceStatus
	var recoveredServices []string
	var emailMessage string

	if err := initializeYamlFiles(configFilename, serviceStatusFilename, &config, &serviceStatus); err != nil {
		fmt.Println(err)
	}

	for _, service := range config.Services {
		isRunning, err := service.isRunning(config.Timeout)

		if err != nil {
			serviceStatus.addAffected(service.Url)
			serviceStatus.incrementDownCount()
		}

		if isRunning && serviceStatus.isAffected(service.Url) {
			serviceStatus.removeAffected(service.Url)
			serviceStatus.decrementDownCount()
			recoveredServices = append(recoveredServices, service.Url)
		}
	}

	if err := serviceStatus.Write(serviceStatusFilename); err != nil {
		fmt.Println(err)
	}

	isDownLimitExceeded := serviceStatus.DownCount >= config.DownLimit
	areServicesRecovered := len(recoveredServices) > 0

	if isDownLimitExceeded || areServicesRecovered {
		emailMessage = getMessage(serviceStatus.AffectedServices, recoveredServices, config.Frequency)
		sendMail(config.MailingList, emailMessage)
		fmt.Println(emailMessage)
	}

}

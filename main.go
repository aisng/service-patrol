package main

import (
	"fmt"
)

const (
	defaultTimeout        uint   = 3
	defaultLimit          uint   = 2
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
			// fmt.Printf("Service %s is down: '%s'\n", service.Url, err)
			serviceStatus.addAffected(service.Url)
			serviceStatus.incrementDownCount()
		}

		if isRunning && serviceStatus.isAffected(service.Url) {
			serviceStatus.removeAffected(service.Url)
			serviceStatus.decrementDownCount()
			recoveredServices = append(recoveredServices, service.Url)

			fmt.Println("Recovered: ", service.Url)
		}

		if isRunning {
			fmt.Printf("Service %s is running\n", service.Url)
		}
	}

	if len(recoveredServices) > 0 {
		fmt.Println(recoveredServices)
		// send recovered email
	}

	isDownLimitExceeded := serviceStatus.DownCount >= config.DownLimit
	areServicesRecovered := len(recoveredServices) > 0

	// TODO: figure out how to avoid two emails being sent at the same time, so that sth recovered and
	// the list of down ones would be in one message
	// solution could be recoveredExists && downExists == true -> send 3rd email template that there are down ones and
	// recovered ones
	emailMessage = getMessage(serviceStatus.AffectedServices, recoveredServices, config.Frequency)

	if isDownLimitExceeded && areServicesRecovered {
		fmt.Println(emailMessage)
		// send services down email + recovered
		sendMail(config.MailingList, emailMessage)
	}

	if isDownLimitExceeded && !areServicesRecovered {
		// send services down
		fmt.Println(emailMessage)
		sendMail(config.MailingList, emailMessage)

	}

	if areServicesRecovered && !isDownLimitExceeded {
		//send services recovered
		fmt.Println(emailMessage)
		sendMail(config.MailingList, emailMessage)

	}

	fmt.Println("-------------------------------------------")
	serviceStatus.Write(serviceStatusFilename)

}

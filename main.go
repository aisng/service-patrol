package main

import (
	"fmt"
)

const (
	defaultTimeout        int    = 3
	defaultLimit          int    = 2
	configFilename        string = "config.yaml"
	serviceStatusFilename string = "service-status.yaml"
)

func main() {

	var config Config
	var serviceStatus ServiceStatus
	var recoveredServices []string

	if err := initializeYamlFiles(configFilename, serviceStatusFilename, &config, &serviceStatus); err != nil {
		fmt.Println(err)
	}

	for _, service := range config.Services {
		isRunning, err := service.isRunning(config.Timeout)

		if err != nil {
			fmt.Printf("Service %s is down: '%s'\n", service.Url, err)
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

	// TODO: figure out how to avoid two emails being sent at the same time, so that sth recovered and
	// the list of down ones
	// solution could be recoveredExists && downExists == true -> send 3rd email template that there are down ones and
	// recovered ones
	if serviceStatus.DownCount >= config.DownLimit {
		fmt.Println(serviceStatus.AffectedServices)
		// send services down email
	}

	fmt.Println("-------------------------------------------")
	serviceStatus.Write(serviceStatusFilename)
	// fmt.Printf("Down count: %d\n", serviceStatus.DownCount)
	// fmt.Println(config.DownLimit > serviceStatus.DownCount)
	// fmt.Printf("RECOVERED: %v", recoveredServices)
	// sendMail(config.MailingList)
}

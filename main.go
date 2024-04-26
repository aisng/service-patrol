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

	if err := initializeYamlFiles(configFilename, serviceStatusFilename, &config, &serviceStatus); err != nil {
		fmt.Println(err)
	}

	services := config.Services

	for _, service := range services {
		isRunning, err := service.isRunning(config.Timeout)

		if err != nil {
			fmt.Printf("Service %s is down: '%s'\n", service.Url, err)
			serviceStatus.addAffected(service.Url)
			serviceStatus.incrementDownCount()
		}

		if isRunning && serviceStatus.isAffected(service.Url) {
			serviceStatus.removeAffected(service.Url)
			serviceStatus.decrementDownCount()
			fmt.Println("Recovered: ", service.Url)
		}

		if isRunning {
			fmt.Printf("Service %s is running\n", service.Url)
		}
	}

	fmt.Println("-------------------------------------------")
	serviceStatus.Write(serviceStatusFilename)
	// fmt.Printf("Down count: %d\n", serviceStatus.DownCount)
	// fmt.Println(config.DownLimit > serviceStatus.DownCount)
	// fmt.Printf("RECOVERED: %v", recoveredServices)
	// sendMail(config.MailingList)
}

// func checkServicesStatus(services []Service, serviceStatus *ServiceStatus) recoveredServices {
// 	recoveredServices := []string{}

// 	for _, service := range services {
// 		isRunning, err := isServiceRunning(service.Url, service.Timeout)
// 		if err != nil {
// 			return err
// 		}
// 		if isRunning {
// 			fmt.Printf("Service %s is running\n", service.Url)
// 			if serviceStatus.IsAffected(service.Url) {
// 				serviceStatus.RemoveAffected(service.Url)
// 				serviceStatus.DecrementDownCount()
// 				recoveredServices = append(recoveredServices, service.Url)
// 				fmt.Println("Recovered: ", service.Url)
// 			}
// 		} else {
// 			serviceStatus.AddAffected(service.Url)
// 			serviceStatus.IncrementDownCount()
// 			fmt.Printf("Service %s is down: '%s'\n", service.Url, err)
// 		}
// 	}
// }

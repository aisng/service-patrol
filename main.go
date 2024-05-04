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
	var serviceStatus ServiceStatus
	var recoveredServices []string
	var affectedServices []string
	var emailMessage string
	var yamlFiles = make(map[string]YamlData)

	// yamlFiles := map[string]YamlData{
	// 	configFilename:        &config,
	// 	serviceStatusFilename: &serviceStatus,
	// }
	yamlFiles[configFilename] = &config
	yamlFiles[serviceStatusFilename] = &serviceStatus

	if err := initializeYamlFiles(yamlFiles); err != nil {
		fmt.Println(err)
	}

	client := http.Client{
		Timeout: time.Second * time.Duration(config.Timeout),
	}

	for _, serviceUrl := range config.Services {
		isRunning, err := isServiceRunning(client, serviceUrl)

		if err != nil {
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
		emailMessage = getMessage(serviceStatus.AffectedServices, recoveredServices, config.Frequency)
		// sendMail(config.MailingList, emailMessage)
		fmt.Println(emailMessage)
	}

}

func isServiceRunning(client http.Client, url string) (bool, error) {
	resp, err := client.Head(url)
	if err != nil {
		return false, err
	}
	resp.Body.Close()
	return true, nil
}

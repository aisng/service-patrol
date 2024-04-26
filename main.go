package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	config := &Config{}
	serviceStatus := &ServiceStatus{}

	if err := initializeYaml(configFilename, serviceStatusFilename, config, serviceStatus); err != nil {
		fmt.Println(err)
	}

	services := config.Services

	for _, service := range services {
		isRunning, err := requestHeadFromService(service.Url, service.Timeout)
		if isRunning {
			fmt.Printf("Service %s is running\n", service.Url)
			if serviceStatus.IsAffected(service.Url) {
				serviceStatus.RemoveAffected(service.Url)
				serviceStatus.DecrementDownCount()
				fmt.Println("Recovered: ", service.Url)
			}

			// TODO: send mail if there was sth down and is up now
		} else {
			serviceStatus.AddAffected(service.Url)
			serviceStatus.IncrementDownCount()
			fmt.Printf("Service %s is down: '%s'\n", service.Url, err)
		}
	}

	fmt.Println("-------------------------------------------")
	serviceStatus.Write(serviceStatusFilename)
	fmt.Printf("Down count: %d\n", serviceStatus.DownCount)
	fmt.Println(config.Limit > serviceStatus.DownCount)
}

func requestHeadFromService(url string, timeout int) (bool, error) {
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	resp, err := client.Head(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}

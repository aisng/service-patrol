package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Services struct {
	Services []struct {
		Url     string `yaml:"url"`
		Timeout int    `yaml:"timeout"`
	} `yaml:"services"`
}

func main() {

	servicesData, err := parseConfigData("config.yaml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(servicesData)
	for _, service := range servicesData.Services {
		isRunning, err := pingService(service.Url, service.Timeout)
		if isRunning {
			fmt.Printf("Service %s is running\n", service.Url)
		} else {
			fmt.Printf("Service %s is down. Reason: %s\n", service.Url, err)
		}

		fmt.Println("-------------------------------------------")

	}

}

func pingService(url string, timeout int) (bool, error) {
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

func parseConfigData(path string) (*Services, error) {
	data, err := os.ReadFile(path)
	fmt.Println(data)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var services Services

	err = yaml.Unmarshal(data, &services)
	fmt.Printf("servicefromparser: %v", services)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &services, nil
}

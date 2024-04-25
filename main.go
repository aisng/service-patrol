package main

import (
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	Url     string `yaml:"url"`
	Timeout int    `yaml:"timeout"`
}

type Config struct {
	Services []Service `yaml:"services"`
	Limit    int       `yaml:"limit"`
}

const defaultTimeout int = 3
const defaultLimit int = 2
const configPath string = "./config.yaml"

func main() {

	var downCount int

	config, err := getConfig(configPath)

	if err != nil {
		fmt.Println(err)
	}

	services := config.Services

	for _, service := range services {
		isRunning, err := requestHeadFromService(service.Url, service.Timeout)
		// TODO: read service-status.json to check what was down before
		if isRunning {
			// TODO: send mail if there was sth down and is up now
			fmt.Printf("Service %s is running\n", service.Url)
		} else {
			fmt.Printf("Service %s is down: '%s'\n", service.Url, err)
			downCount += 1
			// TODO: read service-status.json
			// TODO: count how many down
			// TODO: if more than limit, send email
			// TODO: save in output.txt(?) which one is down
			//
		}
		fmt.Println("-------------------------------------------")
	}
	fmt.Printf("Down count: %d\n", downCount)
	fmt.Println(config.Limit > downCount)

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

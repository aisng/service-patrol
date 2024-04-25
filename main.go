package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// type Services struct {
// 	Services []struct {
// 		Url     string `yaml:"url"`
// 		Timeout int    `yaml:"timeout"`
// 	} `yaml:"services"`
// }

type Service struct {
	Url     string `yaml:"url"`
	Timeout int    `yaml:"timeout"`
}

type Config struct {
	Services []Service `yaml:"services"`
	Limit    int       `yaml:"limit"`
}

const defaultTimeout int = 3
const configPath string = "./config.yaml"

func main() {

	var servicesData Config

	if _, err := os.Stat(configPath); err == nil {

		fmt.Println("'config.yaml' found")
		configData, err := parseConfigData(configPath)

		if err != nil {
			fmt.Println("Parse 'config.yaml' error:", err)
		}

		servicesData = *configData

	} else {
		fmt.Println("config.yaml not found. Generating with default values...")
		configData, err := generateDefaultConfig()

		if err != nil {
			fmt.Println("Generate config.yaml error:", err)
		}
		servicesData = *configData
	}

	for _, service := range servicesData.Services {
		isRunning, err := pingService(service.Url, service.Timeout)

		if isRunning {
			fmt.Printf("Service %s is running\n", service.Url)
		} else {
			fmt.Printf("Service %s is down: '%s'\n", service.Url, err)
			// TODO: count how many down
			// TODO: send email
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

func parseConfigData(path string) (*Config, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	fmt.Printf("servicefromparser: %v\n", config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &config, nil
}

func generateDefaultConfig() (*Config, error) {
	defaultConfig := Config{
		Services: []Service{
			{Url: "http://10.162.222.151/", Timeout: defaultTimeout},
			{Url: "https://prod.alm.gpdm.fresenius.com", Timeout: defaultTimeout},
			{Url: "http://desw-lizenz.schweinfurt.germany.fresenius.de", Timeout: defaultTimeout},
			{Url: "https://central.artifactory.alm.gpdm.fresenius.com", Timeout: defaultTimeout},
			{Url: "https://qdok.ads.fresenius.com/", Timeout: defaultTimeout},
			{Url: "https://www.lrytas.lt", Timeout: defaultTimeout},
		},
		Limit: 2,
	}

	yamlData, err := yaml.Marshal(&defaultConfig)

	if err != nil {
		return nil, err
	}

	// 0644 permission to read, write
	if err = os.WriteFile("config.yaml", yamlData, 0644); err != nil {
		return nil, err
	}
	return &defaultConfig, nil

}

package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func getConfig(path string) (*Config, error) {
	// check if config.yaml exists
	_, err := os.Stat(path)

	if err != nil {
		fmt.Println("'config.yaml' not found. Generating with default values.")
		return generateDefaultConfig()
	}

	fmt.Println("'config.yaml' found.")
	return parseConfigData(path)
}

func parseConfigData(path string) (*Config, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	// fmt.Printf("servicefromparser: %v\n", config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &config, nil
}

func generateDefaultConfig() (*Config, error) {
	defaultConfig := Config{
		Services: []Service{
			{Url: "http://10.162.222.151", Timeout: defaultTimeout},
			{Url: "https://prod.alm.gpdm.fresenius.com", Timeout: defaultTimeout},
			{Url: "http://desw-lizenz.schweinfurt.germany.fresenius.de", Timeout: defaultTimeout},
			{Url: "https://central.artifactory.alm.gpdm.fresenius.com", Timeout: defaultTimeout},
			{Url: "https://qdok.ads.fresenius.com", Timeout: defaultTimeout},
			{Url: "https://www.lrytas.lt", Timeout: defaultTimeout},
		},
		Limit: defaultLimit,
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

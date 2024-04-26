package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type YamlData interface {
	Write(filename string) error
	Read(filename string) error
	GenerateDefault()
}

func writeYaml(filename string, yd YamlData) error {
	yamlData, err := yaml.Marshal(yd)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, yamlData, 0644)
}

func readYaml(filename string, yd YamlData) error {
	yamlData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlData, yd)
}

func initializeYamlFiles(configFilename, serviceStatusFilename string, config *Config, serviceStatus *ServiceStatus) error {
	yamlFiles := map[string]YamlData{
		configFilename:        config,
		serviceStatusFilename: serviceStatus,
	}

	for name, data := range yamlFiles {
		_, err := os.Stat(name)
		if err == nil {
			if err := data.Read(name); err != nil {
				return err
			}
		} else if os.IsNotExist(err) {
			data.GenerateDefault()

			if err := data.Write(name); err != nil {
				return err
			}
		}

	}
	return nil
}

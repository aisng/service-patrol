package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

// TODO: reconsider the necessity of parametrizing interface methods
// when default values can be passed straight to writeYaml and readYaml

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

func initializeYamlFiles(filesMap map[string]YamlData) error {
	for filename, data := range filesMap {
		_, err := os.Stat(filename)

		if err != nil {
			if os.IsNotExist(err) {
				data.GenerateDefault()
				err = data.Write(filename)
				if err != nil {
					return err
				}
			}
			return err
		}

		err = data.Read(filename)

		if err != nil {
			return err
		}
	}
	return nil
}

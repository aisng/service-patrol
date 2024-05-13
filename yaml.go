package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type YamlData interface {
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

	err = yaml.UnmarshalStrict(yamlData, yd)
	if err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			return fmt.Errorf("error reading '%s': file contains incompatible value type:\n %v", filename, err)
		} else if strings.Contains(err.Error(), "not found in type") {
			return fmt.Errorf("error reading '%s': file contains unexpected key(s):\n %v", filename, err)
		} else {
			return err
		}
	}
	return nil
}

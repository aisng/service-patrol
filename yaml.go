package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type YamlData interface{}

func writeYaml(filename string, yd YamlData) error {
	yamlData, err := yaml.Marshal(yd)
	if err != nil {
		return err
	}

	filePath, err := getFilePath(filename)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, yamlData, 0644)
}

func readYaml(filename string, yd YamlData) error {
	filePath, err := getFilePath(filename)
	if err != nil {
		return err
	}

	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.UnmarshalStrict(yamlData, yd)
	if err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			return fmt.Errorf("file contains incompatible value type:\n %v", err)
		} else if strings.Contains(err.Error(), "not found in type") {
			return fmt.Errorf("file contains unexpected key(s):\n %v", err)
		}
		return err
	}
	return nil
}

func getFilePath(filename string) (string, error) {
	appPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get application path: %v", err)
	}
	return filepath.Join(filepath.Dir(appPath), filename), nil
}

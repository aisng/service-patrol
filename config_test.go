package main

import (
	"os"
	"reflect"
	"testing"
)

const testConfigFilename string = "test_config.yaml"

func TestReadAndWriteConfig(t *testing.T) {
	var readConfig Config

	expectedConfig := Config{
		DownLimit:   10,
		Timeout:     9,
		Frequency:   8,
		Services:    []string{"service1", "service2"},
		MailingList: []string{"mail1@example.com", "mail2@example.com"},
	}

	defer os.Remove(testConfigFilename)

	err := expectedConfig.Write(testConfigFilename)
	if err != nil {
		t.Errorf("Error writing Config: %v", err)
	}

	err = readConfig.Read(testConfigFilename)

	if err != nil {
		t.Errorf("Error reading Config: %v", err)
	}

	if !reflect.DeepEqual(readConfig, expectedConfig) {
		t.Errorf("Expected: %v, got: %v", expectedConfig, readConfig)
	}
}

func TestGenerateDefaultConfig(t *testing.T) {
	defaultConfig := &Config{}
	defaultConfig.GenerateDefault()

	expectedConfig := &Config{}
	expectedConfig.GenerateDefault()

	if !reflect.DeepEqual(defaultConfig, expectedConfig) {
		t.Errorf("Expected: %v, got: %v", expectedConfig, defaultConfig)

	}
}

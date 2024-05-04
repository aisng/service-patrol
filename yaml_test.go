package main

import (
	"os"
	"testing"
)

type MockYamlData struct {
	WrittenData           string
	ReadData              string
	GenerateDefaultCalled bool
}

func (m *MockYamlData) Write(filename string) error {
	m.WrittenData = "written content"
	return nil
}

func (m *MockYamlData) Read(filename string) error {
	m.ReadData = "read content"
	return nil
}

func TestWriteYaml(t *testing.T) {
	mockData := &MockYamlData{}

	err := mockData.Write("test_write.yaml")
	if err != nil {
		t.Errorf("error writing yaml: %v", err)
	}

	expected := "written content"

	if mockData.WrittenData != expected {
		t.Errorf("expected content: %v, got: %v", expected, mockData.WrittenData)
	}

	os.Remove("test_write.yaml")
}

func TestReadYaml(t *testing.T) {
	mockData := &MockYamlData{}

	err := mockData.Read("test_read.yaml")
	if err != nil {
		t.Errorf("error writing yaml: %v", err)
	}

	expected := "read content"

	if mockData.ReadData != expected {
		t.Errorf("expected content: %v, got: %v", expected, mockData.ReadData)
	}

	os.Remove("test_read.yaml")
}

func TestInitializeYamlData(t *testing.T) {
	config := &Config{
		DownLimit: 3,
		Timeout:   10,
		Services: []Service{
			{Url: "https://www.example.com"},
		},
		MailingList: []string{"someone@example.com"},
	}

	serviceStatus := &ServiceStatus{}

	configFilename := "config_test.yaml"
	serviceStatusFilename := "service-status_test.yaml"

	defer os.Remove(configFilename)
	defer os.Remove(serviceStatusFilename)

	yamlFilesMap := map[string]YamlData{
		configFilename:        config,
		serviceStatusFilename: serviceStatus,
	}

	err := initializeYamlFiles(yamlFilesMap)

	if err != nil {
		t.Errorf("error initializing yaml files: %v", err)
	}

	if config.DownLimit != defaultLimit {
		t.Errorf("error reading config.DownLimit. Expected: %d, got: %d", defaultLimit, config.DownLimit)
	}

	if config.Timeout != defaultTimeout {
		t.Errorf("error reading config.Timeout. Expected: %d, got: %d", defaultTimeout, config.Timeout)
	}

	if config.Frequency != defaultFrequency {
		t.Errorf("error reading config.Frequency. Expected: %d, got: %d", defaultFrequency, config.Frequency)
	}

	// t.Logf("config %v", config)

	// testcases := []struct {
	// 	config        Config
	// 	serviceStatus ServiceStatus
	// }{
	// 	{
	// 		config: Config{
	// 			DownLimit: 3,
	// 			Timeout:   10,
	// 			Services: []Service{
	// 				{Url: "https://www.example.com"},
	// 			},
	// 			MailingList: []string{"someone@example.com"},
	// 		},
	// 	},
	// 	{
	// 		serviceStatus: ServiceStatus{
	// 			DownCount:        0,
	// 			AffectedServices: []string{},
	// 		},
	// 	},
	// 	{
	// 		config: Config{},
	// 	},
	// 	{
	// 		serviceStatus: ServiceStatus{},
	// 	},
	// }

	// config := &MockYamlData{}
	// serviceStatus := &MockYamlData{}
	// for _, tc := range testcases {

	// 	err := initializeYamlFiles("config_test.yaml", "service_status_test.yaml", &tc.config, &tc.serviceStatus)

	// 	if err != nil {
	// 		t.Errorf("expected %v, got %v", tc.config, err)
	// 	}
	// }
}

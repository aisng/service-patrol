package main

import (
	"errors"
	"io/fs"
	"os"
	"testing"
)

const (
	testFilename string = "test.yaml"
)

type MockYamlData struct {
	Content string
}

func (m *MockYamlData) Write(filename string) error {
	return writeYaml(filename, m)
}

func (m *MockYamlData) Read(filename string) error {
	return readYaml(filename, m)
}

func (m *MockYamlData) GenerateDefault() {
	m.Content = "default"
}

func TestWriteAndReadYaml(t *testing.T) {
	subtests := []struct {
		name            string
		content         string
		expectedContent string
		expectedErr     error
	}{
		{
			name:            "Write",
			content:         "Written content",
			expectedContent: "Written content",
			expectedErr:     nil,
		},
		{
			name:            "ReadFound",
			content:         "Written content",
			expectedContent: "Written content",
			expectedErr:     nil,
		},
		{
			name:            "ReadNotFound",
			content:         "",
			expectedContent: "",
			expectedErr:     fs.ErrNotExist,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			mockData := &MockYamlData{
				Content: subtest.content,
			}

			defer os.Remove(testFilename)
			switch subtest.name {
			case "Write":
				err := mockData.Write(testFilename)
				if err != nil {
					t.Errorf("Error writing yaml: %v", err)
				}
			case "Read":
				err := mockData.Read(testFilename)
				if err != nil {
					t.Errorf("Error reading yaml: %v", err)
				}
			case "ReadNotFound":
				err := mockData.Read("non-existant.yaml")
				if !errors.Is(err, subtest.expectedErr) {
					t.Errorf("Expected: %v, got: %v", subtest.expectedErr, err)
				}
			}

		})
	}

}
func TestWrite(t *testing.T) {
	mockData := &MockYamlData{
		Content: "written",
	}

	defer os.Remove(testFilename)
	err := mockData.Write(testFilename)
	if err != nil {
		t.Errorf("Error writing yaml: %v", err)
	}

	expected := "written"

	if mockData.Content != expected {
		t.Errorf("Expected: %s, got: %s", expected, mockData.Content)
	}

}

func TestRead(t *testing.T) {
	mockData := &MockYamlData{
		Content: "read written",
	}

	defer os.Remove(testFilename)
	mockData.Write(testFilename)

	err := mockData.Read(testFilename)
	if err != nil {
		t.Errorf("Error reading yaml: %v", err)
	}

	expected := "read written"

	if mockData.Content != expected {
		t.Errorf("Expected: %s, got: %s", expected, mockData.Content)
	}

}

func TestInitializeYamlFiles(t *testing.T) {
	defaultData := &MockYamlData{}

	filesMap := map[string]YamlData{testFilename: defaultData}

	defer os.Remove(testFilename)

	err := initializeYamlFiles(filesMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "default"

	if defaultData.Content != "default" {
		t.Errorf("Expected: %s, got: %s", expected, defaultData.Content)
	}

}

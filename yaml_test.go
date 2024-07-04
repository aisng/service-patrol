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

	defer os.Remove(testFilename)

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			mockData := &MockYamlData{
				Content: subtest.content,
			}

			switch subtest.name {
			case "Write":
				err := mockData.Write(testFilename)
				if err != nil {
					t.Errorf("Error writing yaml: %v", err)
				}
				if subtest.content != subtest.expectedContent {
					t.Errorf("Expected: %s, got %s", subtest.expectedContent, subtest.content)
				}
			case "Read":
				err := mockData.Read(testFilename)
				if err != nil {
					t.Errorf("Error reading yaml: %v", err)
				}
				if subtest.content != subtest.expectedContent {
					t.Errorf("Expected: %s, got %s", subtest.expectedContent, subtest.content)
				}
			case "ReadNotFound":
				err := mockData.Read("/random/non-existant.yaml")
				if !errors.Is(err, subtest.expectedErr) {
					t.Errorf("Expected: %v, got: %v", subtest.expectedErr, err)
				}
			}

		})
	}

}

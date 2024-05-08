package main

// import (
// 	"os"
// 	"reflect"
// 	"testing"
// )

// const testConfigFilename string = "test_config.yaml"

// type MockConfig struct {
// 	DownLimit   any
// 	Timeout     any
// 	Frequency   any
// 	Services    any
// 	MailingList any
// }

// func TestConfigRead(t *testing.T) {
// 	subtests := []struct {
// 		name           string
// 		expectedResult MockConfig
// 		expectedErr    error
// 	}{
// 		{
// 			name: "FoundPositive",
// 			expectedResult: MockConfig{
// 				DownLimit:   5,
// 				Timeout:     10,
// 				Frequency:   3,
// 				Services:    []string{"service1", "service2"},
// 				MailingList: []string{"example@example.net"},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			name: "FoundNegativeTypeErr",
// 			expectedResult: MockConfig{
// 				DownLimit:   "5",
// 				Timeout:     10,
// 				Frequency:   3,
// 				Services:    []int{1, 2},
// 				MailingList: []string{"example@example.net"},
// 			},
// 			expectedErr: nil,
// 		},
// 	}

// 	defer os.Remove(testConfigFilename)

// 	err := expectedConfig.Write(testConfigFilename)
// 	if err != nil {
// 		t.Errorf("Error writing Config: %v", err)
// 	}

// 	err = readConfig.Read(testConfigFilename)

// 	if err != nil {
// 		t.Errorf("Error reading Config: %v", err)
// 	}

// 	if !reflect.DeepEqual(readConfig, expectedConfig) {
// 		t.Errorf("Expected: %v, got: %v", expectedConfig, readConfig)
// 	}
// }

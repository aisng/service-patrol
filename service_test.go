package main

import (
	"os"
	"reflect"
	"testing"
)

const testServiceStatusFilename string = "test_service-status.yaml"

func TestServiceStatus(t *testing.T) {
	subtests := []struct {
		name     string
		initial  ServiceStatus
		expected ServiceStatus
	}{
		{
			name: "WriteAndRead",
			initial: ServiceStatus{
				DownCount:        40,
				AffectedServices: []string{"service1", "service2"},
			},
			expected: ServiceStatus{
				DownCount:        0,
				AffectedServices: []string{"service1", "service2"},
			},
		},
		{
			name: "Increment",
			initial: ServiceStatus{
				DownCount:        1,
				AffectedServices: []string{"service1", "service2"},
			},
			expected: ServiceStatus{
				DownCount:        2,
				AffectedServices: []string{"service1", "service2"},
			},
		},
		{
			name: "Decrement",
			initial: ServiceStatus{
				DownCount:        6,
				AffectedServices: []string{"service1", "service2"},
			},
			expected: ServiceStatus{
				DownCount:        5,
				AffectedServices: []string{"service1", "service2"},
			},
		},
		{
			name: "isAffected",
			initial: ServiceStatus{
				DownCount:        0,
				AffectedServices: []string{"service1", "service2"},
			},
			expected: ServiceStatus{
				DownCount:        0,
				AffectedServices: []string{"service1", "service2"},
			},
		},
		{
			name: "isNotAffected",
			initial: ServiceStatus{
				DownCount:        40,
				AffectedServices: []string{"service1", "service2"},
			},
			expected: ServiceStatus{
				DownCount:        5,
				AffectedServices: []string{"service3", "service1"},
			},
		},
	}

	var readServiceStatus ServiceStatus

	defer os.Remove(testServiceStatusFilename)

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			switch subtest.name {
			case "WriteAndRead":
				err := subtest.initial.Write(testServiceStatusFilename)
				if err != nil {
					t.Errorf("Error writing ServiceStatus: %v", err)
				}
				err = readServiceStatus.Read(testServiceStatusFilename)
				if err != nil {
					t.Errorf("Error reading ServiceStatus: %v", err)
				}

				if !reflect.DeepEqual(readServiceStatus, subtest.expected) {
					t.Errorf("Expected: %v, got: %v", subtest.expected, readServiceStatus)
				}

			case "Increment":
				readServiceStatus.DownCount = subtest.initial.DownCount
				readServiceStatus.incrementDownCount()

				if !reflect.DeepEqual(readServiceStatus, subtest.expected) {
					t.Errorf("Expected: %v, got: %v", subtest.expected, readServiceStatus)
				}

			case "Decrement":
				readServiceStatus.DownCount = subtest.initial.DownCount
				readServiceStatus.decrementDownCount()

				if !reflect.DeepEqual(readServiceStatus, subtest.expected) {
					t.Errorf("Expected: %v, got: %v", subtest.expected, readServiceStatus)
				}

			case "isAffected":
				if !readServiceStatus.isAffected(subtest.expected.AffectedServices[0]) {
					t.Errorf("Expected: %v, got: %v", true, false)
				}

			case "isNotAffected":
				if readServiceStatus.isAffected(subtest.expected.AffectedServices[0]) {
					t.Errorf("Expected: %v, got: %v", false, true)
				}
			}
		})
	}
}

// func TestReadAndWriteServiceStatus(t *testing.T) {
// 	var readServiceStatus ServiceStatus

// 	writtenServiceStatus := ServiceStatus{
// 		DownCount:        40,
// 		AffectedServices: []string{"service1", "service2"},
// 	}

// 	expectedServiceStatus := ServiceStatus{
// 		DownCount:        0,
// 		AffectedServices: []string{"service1", "service2"},
// 	}

// 	defer os.Remove(testServiceStatusFilename)

// 	err := writtenServiceStatus.Write(testServiceStatusFilename)
// 	if err != nil {
// 		t.Errorf("Error writing config: %v", err)
// 	}

// 	err = readServiceStatus.Read(testServiceStatusFilename)

// 	if err != nil {
// 		t.Errorf("Error reading config: %v", err)
// 	}

// 	if !reflect.DeepEqual(expectedServiceStatus, readServiceStatus) {
// 		t.Errorf("Expected: %v, got: %v", expectedServiceStatus, readServiceStatus)
// 	}
// }

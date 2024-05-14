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
				DownCount:    40,
				DownServices: []string{"service1", "service2"},
			},
			expected: ServiceStatus{
				DownCount:    0,
				DownServices: []string{"service1", "service2"},
			},
		},
		{
			name: "Increment",
			initial: ServiceStatus{
				DownCount:    1,
				DownServices: []string{"service1", "service2"},
			},
			expected: ServiceStatus{
				DownCount:    2,
				DownServices: []string{"service1", "service2"},
			},
		},
		{
			name: "Decrement",
			initial: ServiceStatus{
				DownCount:    6,
				DownServices: []string{"service1", "service2"},
			},
			expected: ServiceStatus{
				DownCount:    5,
				DownServices: []string{"service1", "service2"},
			},
		},
		{
			name: "IsAffected",
			initial: ServiceStatus{
				DownCount:    0,
				DownServices: []string{"service1", "service2"},
			},
			expected: ServiceStatus{
				DownCount:    0,
				DownServices: []string{"service1", "service2"},
			},
		},
		{
			name: "IsNotAffected",
			initial: ServiceStatus{
				DownCount:    40,
				DownServices: []string{"service1", "service2"},
			},
			expected: ServiceStatus{
				DownCount:    5,
				DownServices: []string{"service3", "service1"},
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

			case "IsAffected":
				if !readServiceStatus.isAffected(subtest.expected.DownServices[0]) {
					t.Errorf("Expected: %v, got: %v", true, false)
				}

			case "IsNotAffected":
				if readServiceStatus.isAffected(subtest.expected.DownServices[0]) {
					t.Errorf("Expected: %v, got: %v", false, true)
				}

			}
		})
	}
}

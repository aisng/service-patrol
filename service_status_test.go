package main

import (
	"os"
	"reflect"
	"testing"
)

const testStatusFilename string = "test-status.yaml"

func TestServiceStatus(t *testing.T) {
	subtests := []struct {
		name     string
		initial  Status
		expected Status
	}{
		{
			name: "WriteAndRead",
			initial: Status{
				DownCount:    40,
				DownServices: []string{"service1", "service2"},
			},
			expected: Status{
				DownCount:    0,
				DownServices: []string{"service1", "service2"},
			},
		},
		{
			name: "Increment",
			initial: Status{
				DownCount:    1,
				DownServices: []string{"service1", "service2"},
			},
			expected: Status{
				DownCount:    2,
				DownServices: []string{"service1", "service2"},
			},
		},
		{
			name: "Decrement",
			initial: Status{
				DownCount:    6,
				DownServices: []string{"service1", "service2"},
			},
			expected: Status{
				DownCount:    5,
				DownServices: []string{"service1", "service2"},
			},
		},
		{
			name: "IsAffected",
			initial: Status{
				DownCount:    0,
				DownServices: []string{"service1", "service2"},
			},
			expected: Status{
				DownCount:    0,
				DownServices: []string{"service1", "service2"},
			},
		},
		{
			name: "IsNotAffected",
			initial: Status{
				DownCount:    40,
				DownServices: []string{"service1", "service2"},
			},
			expected: Status{
				DownCount:    5,
				DownServices: []string{"service3", "service1"},
			},
		},
	}

	var readServiceStatus Status

	defer os.Remove(testStatusFilename)

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			switch subtest.name {
			case "WriteAndRead":
				err := subtest.initial.Write(testStatusFilename)
				if err != nil {
					t.Errorf("Error writing ServiceStatus: %v", err)
				}
				err = readServiceStatus.Read(testStatusFilename)
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

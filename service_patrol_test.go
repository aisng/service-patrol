package main

import (
	"os"
	"reflect"
	"testing"
)

func TestServicePatrol(t *testing.T) {
	subtests := []struct {
		name              string
		config            *Config
		prevStatus        *Status
		expectedDown      []string
		expectedRecovered []string
	}{
		{
			name: "AllDownOverLimit",
			config: &Config{
				DownLimit:   2,
				Timeout:     3,
				Frequency:   4,
				Services:    []string{"http://non-existant1", "http://non-existant2", "http://non-existant3"},
				MailingList: []string{"test@example.org"},
			},
			prevStatus:        NewStatus(),
			expectedDown:      []string{"http://non-existant1", "http://non-existant2", "http://non-existant3"},
			expectedRecovered: nil,
		},
		{
			name: "SomeDownNotOverLimit",
			config: &Config{
				DownLimit:   4,
				Timeout:     3,
				Frequency:   4,
				Services:    []string{"http://lrytas.lt", "http://non-existant2", "http://non-existant3"},
				MailingList: []string{"test@example.org"},
			},
			prevStatus:        NewStatus(),
			expectedDown:      []string{"http://non-existant2", "http://non-existant3"},
			expectedRecovered: nil,
		},
		{
			name: "NoneAreDownAndNoneAreRecovered",
			config: &Config{
				DownLimit:   2,
				Timeout:     3,
				Frequency:   4,
				Services:    []string{"http://lrytas.lt", "http://www.delfi.lt", "http://www.lrt.lt"},
				MailingList: []string{"test@example.org"},
			},
			prevStatus:        NewStatus(),
			expectedDown:      nil,
			expectedRecovered: nil,
		},
	}

	defer os.Remove(testStatusFilename)

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			switch subtest.name {
			case "AllDownOverLimit":
				sp := NewServicePatrol(
					subtest.config,
					subtest.prevStatus,
				)
				down, recovered, _ := sp.Start()
				if !reflect.DeepEqual(down, subtest.expectedDown) {
					t.Errorf("expected: %v, got: %v", subtest.expectedDown, down)
				}
				if !reflect.DeepEqual(recovered, subtest.expectedRecovered) {
					t.Errorf("expected: %v, got: %v", subtest.expectedRecovered, recovered)
				}
			case "SomeDownNotOverLimit":
				sp := NewServicePatrol(
					subtest.config,
					subtest.prevStatus,
				)
				down, recovered, _ := sp.Start()
				if !reflect.DeepEqual(down, subtest.expectedDown) {
					t.Errorf("expected: %v, got: %v", subtest.expectedDown, down)
				}
				if !reflect.DeepEqual(recovered, subtest.expectedRecovered) {
					t.Errorf("expected: %v, got: %v", subtest.expectedRecovered, recovered)
				}
			case "NoneAreDownAndNoneAreRecovered":
				sp := NewServicePatrol(
					subtest.config,
					subtest.prevStatus,
				)
				down, recovered, _ := sp.Start()
				if !reflect.DeepEqual(down, subtest.expectedDown) {
					t.Errorf("expected: %v, got: %v", subtest.expectedDown, down)
				}
				if !reflect.DeepEqual(recovered, subtest.expectedRecovered) {
					t.Errorf("expected: %v, got: %v", subtest.expectedRecovered, recovered)
				}
			}
		})
	}
}

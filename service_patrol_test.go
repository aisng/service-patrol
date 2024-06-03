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
			prevStatus:        NewStatus(0, []string{}),
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
			prevStatus:        NewStatus(0, []string{}),
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
			prevStatus:        NewStatus(0, []string{}),
			expectedDown:      nil,
			expectedRecovered: nil,
		},
		{
			name: "SomeAreDownAndSomeAreRecovered",
			config: &Config{
				DownLimit:   2,
				Timeout:     3,
				Frequency:   4,
				Services:    []string{"http://lrytas.lt", "http://www.delfi.lt", "http://www.lraat.lt"},
				MailingList: []string{"test@example.org"},
			},
			prevStatus:        NewStatus(2, []string{"http://lrytas.lt", "http://www.delfi.lt"}),
			expectedDown:      []string{"http://www.lraat.lt"},
			expectedRecovered: []string{"http://lrytas.lt", "http://www.delfi.lt"},
		},
		{
			name: "NoneDownWithoutHttpPrefix",
			config: &Config{
				DownLimit:   2,
				Timeout:     2,
				Frequency:   1,
				Services:    []string{"google.com", "lrytas.lt", "lrt.lt"},
				MailingList: []string{"test@example.org"},
			},
			prevStatus:        NewStatus(0, []string{}),
			expectedDown:      nil,
			expectedRecovered: nil,
		},
		{
			name: "TwoDownWithoutHttpPrefix",
			config: &Config{
				DownLimit:   2,
				Timeout:     2,
				Frequency:   1,
				Services:    []string{"google.com", "lrasdytas.lt", "lrasdat.lt"},
				MailingList: []string{"test@example.org"},
			},
			prevStatus:        NewStatus(0, []string{}),
			expectedDown:      []string{"lrasdytas.lt", "lrasdat.lt"},
			expectedRecovered: nil,
		},
		{
			name: "AllDownRawIp",
			config: &Config{
				DownLimit:     2,
				Timeout:       2,
				Frequency:     1,
				MaxPacketLoss: 5,
				Services:      []string{"192.0.2.0", "0.42.42.42"},
				MailingList:   []string{"test@example.org"},
			},
			prevStatus:        NewStatus(0, []string{}),
			expectedDown:      []string{"192.0.2.0", "0.42.42.42"},
			expectedRecovered: nil,
		},
		{
			name: "NoneDownRawIp",
			config: &Config{
				DownLimit:     2,
				Timeout:       2,
				Frequency:     1,
				MaxPacketLoss: 5,
				Services:      []string{"8.8.8.8", "131.107.255.255", "4.2.2.2"},
				MailingList:   []string{"test@example.org"},
			},
			prevStatus:        NewStatus(0, []string{}),
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
			case "SomeAreDownAndSomeAreRecovered":
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
			case "NoneDownWithoutHttpPrefix":
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
			case "TwoDownWithoutHttpPrefix":
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
			case "AllDownRawIp":
				sp := NewServicePatrol(
					subtest.config,
					subtest.prevStatus,
				)
				down, recovered, _ := sp.Start()
				isRaw := sp.isRawIpAddress(down[0])

				if isRaw != true {
					t.Errorf("%v is not raw ip", down[0])
				}
				if !reflect.DeepEqual(down, subtest.expectedDown) {
					t.Errorf("expected: %v, got: %v", subtest.expectedDown, down)
				}
				if !reflect.DeepEqual(recovered, subtest.expectedRecovered) {
					t.Errorf("expected: %v, got: %v", subtest.expectedRecovered, recovered)
				}
				// case "NoneDownRawIp":
				// 	sp := NewServicePatrol(
				// 		subtest.config,
				// 		subtest.prevStatus,
				// 	)
				// 	down, recovered, _ := sp.Start()

				// 	if !reflect.DeepEqual(down, subtest.expectedDown) {
				// 		t.Errorf("expected: %v, got: %v", subtest.expectedDown, down)
				// 	}
				// 	if !reflect.DeepEqual(recovered, subtest.expectedRecovered) {
				// 		t.Errorf("expected: %v, got: %v", subtest.expectedRecovered, recovered)
				// 	}
			}
		})
	}
}

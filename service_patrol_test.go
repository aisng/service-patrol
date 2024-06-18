package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestServicePatrol(t *testing.T) {
	subtests := []struct {
		name           string
		config         *Config
		status         *Status
		expectedResult [][]string
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
			status:         NewStatus(0, []string{}),
			expectedResult: [][]string{{"http://non-existant1", "http://non-existant2", "http://non-existant3"}, nil},
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
			status:         NewStatus(0, []string{}),
			expectedResult: [][]string{{"http://non-existant2", "http://non-existant3"}, nil},
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
			status:         NewStatus(0, []string{}),
			expectedResult: [][]string{nil, nil},
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
			status:         NewStatus(2, []string{"http://lrytas.lt", "http://www.delfi.lt"}),
			expectedResult: [][]string{{"http://www.lraat.lt"}, {"http://lrytas.lt", "http://www.delfi.lt"}},
		},
		{
			name: "NoneDownWithoutHttpPrefix",
			config: &Config{
				DownLimit:   2,
				Timeout:     5,
				Frequency:   1,
				Services:    []string{"google.com", "lrytas.lt", "lrt.lt"},
				MailingList: []string{"test@example.org"},
			},
			status:         NewStatus(0, []string{}),
			expectedResult: [][]string{nil, nil},
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
			status:         NewStatus(0, []string{}),
			expectedResult: [][]string{{"lrasdytas.lt", "lrasdat.lt"}, nil},
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
			status:         NewStatus(0, []string{}),
			expectedResult: [][]string{{"192.0.2.0", "0.42.42.42"}, nil},
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
			status:         NewStatus(0, []string{}),
			expectedResult: [][]string{nil, nil},
		},
	}

	defer os.Remove(testStatusFilename)

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			switch subtest.name {
			case "AllDownOverLimit":
				sp := NewServicePatrol(
					subtest.config,
					subtest.status,
				)
				down, recovered, _ := sp.Start()
				fmt.Printf("---- TYPE RECOVERED: %v, TYPE EXPECTED RECOVERED: %v\n", reflect.TypeOf(recovered), reflect.TypeOf(subtest.expectedResult[1]))
				fmt.Printf("---- len RECOVERED: %v, len EXPECTED RECOVERED: %v\n", len(recovered), len(subtest.expectedResult[1]))
				if !reflect.DeepEqual(down, subtest.expectedResult[0]) {
					t.Errorf("expected: %v, got: %v", subtest.expectedResult[0], down)
				}
				if !reflect.DeepEqual(recovered, subtest.expectedResult[1]) {
					t.Errorf("expected: %v, got: %v", subtest.expectedResult[1], recovered)
				}
			case "SomeDownNotOverLimit":
				sp := NewServicePatrol(
					subtest.config,
					subtest.status,
				)
				down, _, _ := sp.Start()
				if !reflect.DeepEqual(down, subtest.expectedResult[0]) {
					t.Errorf("expected: %v, got: %v", subtest.expectedResult[0], down)
				}

			case "NoneAreDownAndNoneAreRecovered":
				sp := NewServicePatrol(
					subtest.config,
					subtest.status,
				)
				down, recovered, _ := sp.Start()
				if !reflect.DeepEqual(down, subtest.expectedResult[0]) {
					t.Errorf("expected: %v, got: %v", subtest.expectedResult[0], down)
				}
				if !reflect.DeepEqual(recovered, subtest.expectedResult[1]) {
					t.Errorf("expected: %v, got: %v", subtest.expectedResult[1], recovered)
				}
			case "SomeAreDownAndSomeAreRecovered":
				sp := NewServicePatrol(
					subtest.config,
					subtest.status,
				)
				down, recovered, _ := sp.Start()
				if !reflect.DeepEqual(down, subtest.expectedResult[0]) {
					t.Errorf("expected: %v, got: %v", subtest.expectedResult[0], down)
				}
				if !reflect.DeepEqual(recovered, subtest.expectedResult[1]) {
					t.Errorf("expected: %v, got: %v", subtest.expectedResult[1], recovered)
				}
			case "NoneDownWithoutHttpPrefix":
				sp := NewServicePatrol(
					subtest.config,
					subtest.status,
				)
				down, recovered, _ := sp.Start()
				if !reflect.DeepEqual(down, subtest.expectedResult[0]) {
					t.Errorf("expected: %v, got: %v", subtest.expectedResult[0], down)
				}
				if !reflect.DeepEqual(recovered, subtest.expectedResult[1]) {
					t.Errorf("expected: %v, got: %v", subtest.expectedResult[1], recovered)
				}
			case "TwoDownWithoutHttpPrefix":
				sp := NewServicePatrol(
					subtest.config,
					subtest.status,
				)
				down, recovered, _ := sp.Start()
				if !reflect.DeepEqual(down, subtest.expectedResult[0]) {
					t.Errorf("expected: %v, got: %v", subtest.expectedResult[0], down)
				}
				if !reflect.DeepEqual(recovered, subtest.expectedResult[1]) {
					t.Errorf("expected: %v, got: %v", subtest.expectedResult[1], recovered)
				}
			case "AllDownRawIp":
				sp := NewServicePatrol(
					subtest.config,
					subtest.status,
				)
				down, recovered, err := sp.Start()
				if err != nil {
					t.Skipf("skipping test because of lack of permissions: %v", err)
				} else {
					t.Errorf("Unexpected error: %v", err)
				}

				isRaw := sp.isRawIpAddress(down[0])
				if !isRaw {
					t.Errorf("%v is not raw ip", down[0])
				}
				if !reflect.DeepEqual(down, subtest.expectedResult[0]) {
					t.Errorf("expected: %v, got: %v", subtest.expectedResult[0], down)
				}
				if !reflect.DeepEqual(recovered, subtest.expectedResult[1]) {
					t.Errorf("expected: %v, got: %v", subtest.expectedResult[1], recovered)
				}

			}
		})
	}
}

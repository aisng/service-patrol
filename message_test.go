package main

import (
	"testing"
)

func TestFormatServicesListStr(t *testing.T) {
	services := []string{"service1", "service2"}
	expected := " - service1\n - service2\n"
	result := formatServicesListStr(services)

	if expected != result {
		t.Errorf("Expected: %v, got: %v", expected, result)
	}
}

func TestParseTemplate(t *testing.T) {
	subtests := []struct {
		name     string
		msg      *Message
		expected string
	}{
		{
			name: "ServicesAreDown",
			msg:  NewMessage([]string{"service1", "service2"}, []string{}, 4),
			expected: `Subject: Connection to FMC services lost
Hello,

connection to the pages/IPs below was lost:
 - service1
 - service2

Next check will be made after 4 hours.`,
		},
		{
			name: "ServicesAreRecovered",
			msg:  NewMessage([]string{}, []string{"service1", "service2"}, 5),
			expected: `Subject: Connection to FMC services recovered
Hello,

connection to the pages/IPs below was recovered:
 - service1
 - service2

Next check will be made after 5 hours.`,
		},
		{
			name: "ServicesAreRecoveredAndSomeAreDown",
			msg:  NewMessage([]string{"service3", "service4"}, []string{"service1", "service2"}, 10),
			expected: `Subject: Connection to some FMC services recovered
Hello,

connection to the pages/IPs below was recovered:
 - service1
 - service2

The following pages are still down:
 - service3
 - service4

Next check will be made after 10 hours.`,
		},
		{
			name:     "ServicesAreNotDownAndNotRecovered",
			msg:      NewMessage([]string{}, []string{}, 2),
			expected: "",
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			msg, err := ParseTemplate(subtest.msg, messageTemplate)

			switch subtest.name {
			case "ServicesAreDown":
				if msg != subtest.expected {
					t.Errorf("expected: %s, got %s", subtest.expected, msg)
				}
			case "ServicesAreRecovered":
				if msg != subtest.expected {
					t.Errorf("expected: %s, got %s", subtest.expected, msg)
				}
			case "ServicesAreRecoveredAndSomeAreDown":
				if msg != subtest.expected {
					t.Errorf("expected: %s, got %s", subtest.expected, msg)
				}
			case "ServicesAreNotDownAndNotRecovered":
				if msg != subtest.expected && err.Error() != "struct is empty: nothing to parse" {
					t.Errorf("expected: %s, got %s", subtest.expected, msg)
				}
			}

		},
		)
	}
}

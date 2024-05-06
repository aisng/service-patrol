package main

import (
	"testing"
)

func TestGenerateServicesList(t *testing.T) {
	services := []string{"service1", "service2"}
	expected := " - service1\n - service2\n"
	result := generateServicesList(services)

	if expected != result {
		t.Errorf("Expected: %v, got: %v", expected, result)
	}
}

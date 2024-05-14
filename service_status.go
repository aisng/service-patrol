package main

import (
	"fmt"
	"os"
)

const serviceStatusFilename string = "service-status.yaml"

type ServiceStatus struct {
	DownCount    uint     `yaml:"down_count"`
	DownServices []string `yaml:"down_services"`
}

func NewServiceStatus() *ServiceStatus {
	return &ServiceStatus{
		DownCount:    0,
		DownServices: []string{},
	}
}

func (ss *ServiceStatus) Write(filename string) error {
	return writeYaml(filename, ss)
}

func (ss *ServiceStatus) Read(filename string) error {
	if err := readYaml(filename, ss); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("'%s' not found and will be created\n", serviceStatusFilename)
			NewServiceStatus()
			return nil
		}
	}
	ss.DownCount = 0 // reset the down count each time the app is run
	return nil
}

func (ss *ServiceStatus) incrementDownCount() {
	ss.DownCount++
}

func (ss *ServiceStatus) decrementDownCount() {
	if ss.DownCount > 0 {
		ss.DownCount--
	}
}

func (ss *ServiceStatus) isAffected(url string) bool {
	for _, affectedService := range ss.DownServices {
		if url == affectedService {
			return true
		}
	}
	return false
}

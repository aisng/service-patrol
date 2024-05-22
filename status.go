package main

import (
	"log"
	"os"
)

const statusFilename string = "status.yaml"

type Status struct {
	DownCount    uint     `yaml:"down_count"`
	DownServices []string `yaml:"down_services"`
}

func NewStatus() *Status {
	return &Status{
		DownCount:    0,
		DownServices: []string{},
	}
}

func (ss *Status) Write(filename string) error {
	return writeYaml(filename, ss)
}

func (ss *Status) Read(filename string) error {
	if err := readYaml(filename, ss); err != nil {
		if os.IsNotExist(err) {
			log.Printf("'%s' not found and will be created\n", statusFilename)
			NewStatus()
			return nil
		} else {
			return err
		}
	}
	ss.DownCount = 0 // reset the down count each time the app is run
	return nil
}

func (ss *Status) incrementDownCount() {
	ss.DownCount++
}

func (ss *Status) decrementDownCount() {
	if ss.DownCount > 0 {
		ss.DownCount--
	}
}

func (ss *Status) isAffected(url string) bool {
	for _, affectedService := range ss.DownServices {
		if url == affectedService {
			return true
		}
	}
	return false
}

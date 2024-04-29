package main

import (
	"net/http"
	"time"
)

type Service struct {
	Url string `yaml:"url"`
}

func (s *Service) isRunning(timeout uint) (bool, error) {
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	resp, err := client.Head(s.Url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}

type ServiceStatus struct {
	DownCount        uint     `yaml:"down_count"`
	AffectedServices []string `yaml:"affected_services"`
}

func (ss *ServiceStatus) Write(filename string) error {
	return writeYaml(filename, ss)
}

func (ss *ServiceStatus) Read(filename string) error {
	if err := readYaml(filename, ss); err != nil {
		return err
	}
	ss.DownCount = 0 // reset the down count each time app is run
	return nil
}

func (ss *ServiceStatus) GenerateDefault() {
	ss.DownCount = 0
	ss.AffectedServices = []string{}
}

func (ss *ServiceStatus) addAffected(url string) {
	for _, affectedService := range ss.AffectedServices {
		if affectedService == url {
			return
		}
	}
	ss.AffectedServices = append(ss.AffectedServices, url)
}

func (ss *ServiceStatus) removeAffected(url string) {
	for i, affectedService := range ss.AffectedServices {
		if affectedService == url {
			ss.AffectedServices = append(ss.AffectedServices[:i], ss.AffectedServices[i+1:]...)
			return
		}
	}
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
	for _, affectedService := range ss.AffectedServices {
		if url == affectedService {
			return true
		}
	}
	return false
}

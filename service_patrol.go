package main

import (
	"fmt"
	"log"
	"net/http"
)

type ServicePatrol struct {
	Config            *Config
	PrevStatus        *Status
	Client            *http.Client
	RecoveredServices []string
	DownServices      []string
}

func NewServicePatrol(config *Config, prevStatus *Status) *ServicePatrol {
	return &ServicePatrol{
		Config:     config,
		PrevStatus: prevStatus,
		Client:     NewHttpClient(config.Timeout),
	}
}

func (sp *ServicePatrol) Start() ([]string, []string, error) {
	for _, serviceUrl := range sp.Config.Services {
		isRunning, err := sp.isServiceRunning(serviceUrl)

		if err != nil {
			log.Printf("service down: %s\n", serviceUrl)
		}

		if !isRunning {
			sp.DownServices = append(sp.DownServices, serviceUrl)
			sp.PrevStatus.incrementDownCount()
		}

		if isRunning && sp.PrevStatus.isAffected(serviceUrl) {
			sp.RecoveredServices = append(sp.RecoveredServices, serviceUrl)
			sp.PrevStatus.decrementDownCount()
		}
	}

	// assign found down services to Status struct and write to .yaml
	sp.PrevStatus.DownServices = sp.DownServices
	if err := sp.PrevStatus.Write(statusFilename); err != nil {
		return nil, nil, fmt.Errorf("error writing to %q: %v", statusFilename, err)
	}

	return sp.DownServices, sp.RecoveredServices, nil
}

func (sp *ServicePatrol) IsDownLimitExceeded() bool {
	return sp.PrevStatus.DownCount >= sp.Config.DownLimit
}

func (sp *ServicePatrol) IsRecoveredFound() bool {
	return len(sp.RecoveredServices) > 0
}

func (sp *ServicePatrol) IsDownFound() bool {
	return len(sp.DownServices) > 0
}

func (sp *ServicePatrol) isServiceRunning(url string) (bool, error) {
	resp, err := sp.Client.Head(url)
	if err != nil {
		return false, err
	}
	resp.Body.Close()
	return true, nil
}

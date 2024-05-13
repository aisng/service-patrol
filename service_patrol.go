package main

import (
	"fmt"
	"net/http"
)

type ServicePatrol struct {
	Config            *Config
	PrevStatus        *ServiceStatus
	Client            *http.Client
	RecoveredServices []string
	DownServices      []string
}

func NewServicePatrol(config *Config, prevStatus *ServiceStatus) *ServicePatrol {
	return &ServicePatrol{
		Config:     config,
		PrevStatus: prevStatus,
		Client:     NewClient(config.Timeout),
	}
}

func (sp *ServicePatrol) Start() (down, recovered []string) {
	for _, serviceUrl := range sp.Config.Services {
		isRunning, err := sp.isServiceRunning(serviceUrl)

		if err != nil {
			fmt.Println(err)
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

	sp.PrevStatus.DownServices = sp.DownServices

	if err := sp.PrevStatus.Write(); err != nil {
		fmt.Println(err)
	}

	if sp.isDownLimitExceeded() || sp.isRecoveredFound() {
		return sp.DownServices, sp.RecoveredServices
	}

	return nil, nil
}

func (sp *ServicePatrol) isDownLimitExceeded() bool {
	return sp.PrevStatus.DownCount >= sp.Config.DownLimit
}

// func (sp *ServicePatrol) isDownFound() bool {
// 	return len(sp.DownServices) > 0
// }

func (sp *ServicePatrol) isRecoveredFound() bool {
	return len(sp.RecoveredServices) > 0
}

func (sp *ServicePatrol) isServiceRunning(url string) (bool, error) {
	resp, err := sp.Client.Head(url)
	if err != nil {
		return false, err
	}
	resp.Body.Close()
	return true, nil
}

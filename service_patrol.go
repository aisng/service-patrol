package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	probing "github.com/prometheus-community/pro-bing"
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
	for _, serviceAddress := range sp.Config.Services {
		isRunning, err := sp.isServiceRunning(serviceAddress)

		if err != nil {
			isNotPermitted := strings.Contains(err.Error(), "operation not permitted")
			isHostNotFound := strings.Contains(err.Error(), "no such host")
			isTimeoutExceeded := strings.Contains(err.Error(), "context deadline exceeded")
			isPacketLimitExceeded := strings.Contains(err.Error(), "packet loss limit exceeded")

			if isNotPermitted {
				return nil, nil, fmt.Errorf("cannot ping %q: %v", serviceAddress, err)
			} else if isTimeoutExceeded || isHostNotFound || isPacketLimitExceeded {
				log.Printf("service down: %q: %v", serviceAddress, err)
			} else {
				return nil, nil, err
			}
		}

		if !isRunning {
			sp.DownServices = append(sp.DownServices, serviceAddress)
			sp.PrevStatus.incrementDownCount()
		}

		if isRunning && sp.PrevStatus.isAffected(serviceAddress) {
			sp.RecoveredServices = append(sp.RecoveredServices, serviceAddress)
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

func (sp *ServicePatrol) isServiceRunning(addr string) (bool, error) {
	if sp.isRawIpAddress(addr) {
		stats, err := sp.getPingerStats(addr)
		if err != nil {
			return false, err
		}

		// packet loss amount is in percentage
		if stats.PacketLoss > float64(sp.Config.MaxPacketLoss) {
			return false, fmt.Errorf("packet loss limit exceeded: addr: %s sent: %v received: %v loss: %v",
				addr, stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		}

		return true, nil
	}

	err := sp.sendHeadRequest(addr)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (sp *ServicePatrol) sendHeadRequest(addr string) error {
	hasHttp := strings.HasPrefix(addr, "http://")
	hasHttps := strings.HasPrefix(addr, "https://")

	if !(hasHttp || hasHttps) {
		addr = "http://" + addr
	}

	resp, err := sp.Client.Head(addr)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil

}

func (sp *ServicePatrol) isRawIpAddress(addr string) bool {
	ipv4Regex := regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`)
	return ipv4Regex.MatchString(addr)
}

func (sp *ServicePatrol) getPingerStats(addr string) (*probing.Statistics, error) {
	pinger, err := probing.NewPinger(addr)
	pinger.SetPrivileged(true)

	if err != nil {
		return nil, err
	}

	pinger.Count = 3
	pinger.Timeout = time.Second * time.Duration(sp.Config.Timeout)

	err = pinger.Run()
	if err != nil {
		return nil, err
	}
	return pinger.Statistics(), nil

}

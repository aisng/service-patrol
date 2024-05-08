package main

type ServiceStatus struct {
	DownCount        uint     `yaml:"down_count"`
	AffectedServices []string `yaml:"affected_services"`
}

func NewServiceStatus() *ServiceStatus {
	return &ServiceStatus{
		DownCount:        0,
		AffectedServices: []string{},
	}
}

func (ss *ServiceStatus) Write(filename string) error {
	return writeYaml(filename, ss)
}

func (ss *ServiceStatus) Read(filename string) error {
	if err := readYaml(filename, ss); err != nil {
		return err
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
	for _, affectedService := range ss.AffectedServices {
		if url == affectedService {
			return true
		}
	}
	return false
}

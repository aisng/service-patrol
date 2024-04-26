package main

const (
	serviceStatusFilename string = "service-status.yaml"
)

type Service struct {
	Url     string `yaml:"url"`
	Timeout int    `yaml:"timeout"`
}

type ServiceStatus struct {
	DownCount        int      `yaml:"down_count"`
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

func (ss *ServiceStatus) AddAffected(url string) {
	for _, affectedService := range ss.AffectedServices {
		if affectedService == url {
			return
		}
	}
	ss.AffectedServices = append(ss.AffectedServices, url)
}

func (ss *ServiceStatus) RemoveAffected(url string) {
	for i, affectedService := range ss.AffectedServices {
		if affectedService == url {
			ss.AffectedServices = append(ss.AffectedServices[:i], ss.AffectedServices[i+1:]...)
			return
		}
	}
}

func (ss *ServiceStatus) IncrementDownCount() {
	ss.DownCount++
}

func (ss *ServiceStatus) DecrementDownCount() {
	if ss.DownCount > 0 {
		ss.DownCount--
	}
}

func (ss *ServiceStatus) IsAffected(url string) bool {
	for _, affectedService := range ss.AffectedServices {
		if url == affectedService {
			return true
		}
	}
	return false
}

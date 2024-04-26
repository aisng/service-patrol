package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type YamlData interface {
	Write(filename string) error
	Read(filename string) error
	GenerateDefault()
}

type Service struct {
	Url     string `yaml:"url"`
	Timeout int    `yaml:"timeout"`
}

type Config struct {
	Services []Service `yaml:"services"`
	Limit    int       `yaml:"limit"`
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

func (c *Config) Write(filename string) error {
	return writeYaml(filename, c)
}

func (c *Config) Read(filename string) error {
	return readYaml(filename, c)
}

func (c *Config) GenerateDefault() {
	c.Services = []Service{
		{Url: "http://10.162.222.151", Timeout: defaultTimeout},
		{Url: "https://prod.alm.gpdm.fresenius.com", Timeout: defaultTimeout},
		{Url: "http://desw-lizenz.schweinfurt.germany.fresenius.de", Timeout: defaultTimeout},
		{Url: "https://central.artifactory.alm.gpdm.fresenius.com", Timeout: defaultTimeout},
		{Url: "https://qdok.ads.fresenius.com", Timeout: defaultTimeout},
		{Url: "https://www.lrytas.lt", Timeout: defaultTimeout},
	}
	c.Limit = defaultLimit
}

const defaultTimeout int = 3
const defaultLimit int = 2
const configFilename string = "config.yaml"
const serviceStatusFilename string = "service-status.yaml"

func main() {

	config := &Config{}
	serviceStatus := &ServiceStatus{}

	if err := initializeYaml(configFilename, serviceStatusFilename, config, serviceStatus); err != nil {
		fmt.Println(err)
	}

	services := config.Services

	for _, service := range services {
		isRunning, err := requestHeadFromService(service.Url, service.Timeout)
		if isRunning {
			fmt.Printf("Service %s is running\n", service.Url)
			if serviceStatus.IsAffected(service.Url) {
				serviceStatus.RemoveAffected(service.Url)
				serviceStatus.DecrementDownCount()
				fmt.Println("Recovered: ", service.Url)
			}

			// TODO: send mail if there was sth down and is up now
		} else {
			serviceStatus.AddAffected(service.Url)
			serviceStatus.IncrementDownCount()
			fmt.Printf("Service %s is down: '%s'\n", service.Url, err)
		}
	}

	fmt.Println("-------------------------------------------")
	serviceStatus.Write(serviceStatusFilename)
	fmt.Printf("Down count: %d\n", serviceStatus.DownCount)
	fmt.Println(config.Limit > serviceStatus.DownCount)
}

func requestHeadFromService(url string, timeout int) (bool, error) {
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	resp, err := client.Head(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}

func writeYaml(filename string, yd YamlData) error {
	yamlData, err := yaml.Marshal(yd)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, yamlData, 0644)
}

func readYaml(filename string, yd YamlData) error {
	yamlData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlData, yd)
}

func initializeYaml(configFilename, serviceStatusFilename string, config *Config, serviceStatus *ServiceStatus) error {
	yamlFiles := map[string]YamlData{
		configFilename:        config,
		serviceStatusFilename: serviceStatus,
	}

	for name, data := range yamlFiles {
		_, err := os.Stat(name)
		if err == nil {
			if err := data.Read(name); err != nil {
				return err
			}
		} else if os.IsNotExist(err) {
			data.GenerateDefault()

			if err := data.Write(name); err != nil {
				return err
			}
		}

	}
	return nil
}

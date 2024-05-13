package main

const configFilename string = "config.yaml"

type Config struct {
	DownLimit   uint     `yaml:"down_limit"`
	Timeout     uint     `yaml:"timeout_s"`
	Frequency   uint     `yaml:"frequency_h"`
	Services    []string `yaml:"services"`
	MailingList []string `yaml:"mailing_list"`
}

func (c *Config) Read() error {
	return readYaml(configFilename, c)
}

package main

const configFilename string = "config.yaml"

type Config struct {
	DownLimit     uint     `yaml:"down_limit"`
	MaxPacketLoss uint     `yaml:"packet_loss_limit_percent"`
	Timeout       uint     `yaml:"timeout_s"`
	Frequency     uint     `yaml:"frequency_h"`
	Provider      string   `yaml:"provider"`
	Services      []string `yaml:"services"`
	MailingList   []string `yaml:"mailing_list"`
}

func (c *Config) Read(filename string) error {
	return readYaml(filename, c)
}

package main

type Config struct {
	DownLimit   uint     `yaml:"down_limit"`
	Timeout     uint     `yaml:"timeout_s"`
	Frequency   uint     `yaml:"frequency_h"`
	Services    []string `yaml:"services"`
	MailingList []string `yaml:"mailing_list"`
}

func (c *Config) Read(filename string) error {
	// if err := validateFields(filename, c); err != nil {
	// 	return err
	// }
	return readYaml(filename, c)
}

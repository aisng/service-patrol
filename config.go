package main

type Config struct {
	DownLimit   int       `yaml:"down_limit"`
	Timeout     int       `yaml:"timeout"`
	Services    []Service `yaml:"services"`
	MailingList []string  `yaml:"mailing_list"`
}

func (c *Config) Write(filename string) error {
	return writeYaml(filename, c)
}

func (c *Config) Read(filename string) error {
	return readYaml(filename, c)
}

func (c *Config) GenerateDefault() {
	c.DownLimit = defaultLimit
	c.Timeout = defaultTimeout
	c.Services = []Service{
		{Url: "http://10.162.222.151"},
		{Url: "https://prod.alm.gpdm.fresenius.com"},
		{Url: "http://desw-lizenz.schweinfurt.germany.fresenius.de"},
		{Url: "https://central.artifactory.alm.gpdm.fresenius.com"},
		{Url: "https://qdok.ads.fresenius.com"},
		{Url: "https://www.lrytas.lt"},
	}
	c.MailingList = []string{
		"mordechai.lashaun@dockerbike.com",
		"pcjhwssca2@zlorkun.com",
	}

}

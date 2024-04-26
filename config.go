package main

const ()

type Config struct {
	Services []Service `yaml:"services"`
	Limit    int       `yaml:"limit"`
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

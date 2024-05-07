package main

type Config struct {
	DownLimit   uint     `yaml:"down_limit"`
	Timeout     uint     `yaml:"timeout_s"`
	Frequency   uint     `yaml:"frequency_h"`
	Services    []string `yaml:"services"`
	MailingList []string `yaml:"mailing_list"`
}

// func (c *Config) Read(filename string) error {
// 	return ReadYaml(filename, c)
// }

func (c *Config) Read(filename string) error {
	if err := validateFields(filename, c); err != nil {
		return err
	}
	return readYaml(filename, c)
}

// func (c *Config) Read(filename string) []error {
// 	data := make(map[string]interface{})

// 	if err := ReadYaml(filename, &data); err != nil {
// 		return []error{err}
// 	}

// 	unrecognizedFields := []error{}

// 	for key := range data {
// 		// Check if the field exists in the Config struct
// 		if _, ok := reflect.TypeOf(c).Elem().FieldByName(key); !ok {
// 			erri := fmt.Errorf("unrecognized field '%s' in YAML", key)
// 			unrecognizedFields = append(unrecognizedFields, erri)
// 		}
// 	}
// 	if len(unrecognizedFields) != 0 {
// 		return unrecognizedFields
// 	}
// 	// Unmarshal the YAML data into the Config struct
// 	// if err := ReadYaml(filename, c); err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }

// func (c *Config) GenerateDefault() {
// 	c.DownLimit = defaultLimit
// 	c.Timeout = defaultTimeout
// 	c.Frequency = defaultFrequency
// 	c.Services = []string{
// 		"http://10.162.222.151",
// 		"https://prod.alm.gpdm.fresenius.com",
// 		"http://desw-lizenz.schweinfurt.germany.fresenius.de",
// 		"https://central.artifactory.alm.gpdm.fresenius.com",
// 		"https://qdok.ads.fresenius.com",
// 		"https://www.lrytas.lt",
// 	}
// 	c.MailingList = []string{
// 		"mohammad.abshir@dockerbike.com",
// 	}
// }

// func (c *Config) GenerateDefault() *Config {
// 	return &Config{
// 		DownLimit: defaultLimit,
// 		Timeout:   defaultTimeout,
// 		Frequency: defaultFrequency,
// 		Services: []string{
// 			"http://10.162.222.151",
// 			"https://prod.alm.gpdm.fresenius.com",
// 			"http://desw-lizenz.schweinfurt.germany.fresenius.de",
// 			"https://central.artifactory.alm.gpdm.fresenius.com",
// 			"https://qdok.ads.fresenius.com",
// 			"https://www.lrytas.lt",
// 		},
// 		MailingList: []string{
// 			"mohammad.abshir@dockerbike.com",
// 		},
// 	}
// }

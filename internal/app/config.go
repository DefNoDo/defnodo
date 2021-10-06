package app

type Config struct {
	configPath string
}

func LoadConfig(location string) (config *Config, err error) {
	config = &Config{
		configPath: location,
	}
	return
}

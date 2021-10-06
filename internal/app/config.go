package app

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DataDirectory string `yaml:"data-directory"`
}

func LoadConfig(location string) (config *Config, err error) {
	log.Printf("Loading config from %s\n", location)
	config = &Config{}
	if _, err := os.Stat(location); os.IsNotExist(err) {
		config = loadDefaults(config)
	} else {
		config, err = loadConfig(config, location)
	}
	return
}

func loadConfig(c *Config, location string) (config *Config, err error) {
	file, err := os.Open(location)
	if err != nil {
		log.Printf("Error loading config file '%s'\n", location)
		return
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		log.Print("Error parsing config file\n")
	}

	config = loadDefaults(c)
	return
}

func loadDefaults(c *Config) (config *Config) {
	if c.DataDirectory == "" {
		c.DataDirectory = "~/.defnodo"
	}

	config = c

	return
}

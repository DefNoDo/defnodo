package app

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config is the runtime configuration
type Config struct {
	DataDirectory    string   `yaml:"data-directory"`
	VolumeMounts     []string `yaml:"volume-mounts"`
	DockerDaemonJson string   `yaml:"docker-daemon.json"`
	VM               VMConfig `yaml:"vm"`
	IsService        bool
	Interactive      bool
}

type VMConfig struct {
	Memory   int    `yaml:"memory"`
	Cpus     int    `yaml:"cpus"`
	DiskSize string `yaml:"disk-size"`
}

// Load a configuration from a given location.  Any unset values will
// use the default values associated.  If location does not exist, only
// default values will be used
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

// Load a yaml configuration from disk and then apply default values for
// any unset values
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

// Assign default values for any unset values in c
func loadDefaults(c *Config) (config *Config) {
	if c.DataDirectory == "" {
		c.DataDirectory = "~/.defnodo"
	}
	if len(c.VolumeMounts) == 0 {
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		c.VolumeMounts = []string{homedir}
	}
	c.VM = loadVMDefaults(c.VM)

	config = c

	return
}

func loadVMDefaults(c VMConfig) (config VMConfig) {
	if c.Memory <= 0 {
		c.Memory = 2048
	}
	if c.Cpus <= 0 {
		c.Cpus = 1
	}
	if len(c.DiskSize) == 0 {
		c.DiskSize = "10G"
	}
	config = c
	return
}

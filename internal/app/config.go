package app

import (
	"log"
	"os"
	"path/filepath"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

// Config is the runtime configuration
type Config struct {
	DataDirectory       string   `yaml:"data-directory"`
	VolumeMounts        []string `yaml:"volume-mounts"`
	DockerDaemonJson    string   `default:"" yaml:"docker-daemon.json"`
	VM                  VMConfig `yaml:"vm"`
	IsService           bool     `default:"false"`
	Interactive         bool     `default:"false"`
	ConfigBaseDirectory string
}

type VMConfig struct {
	Memory   int    `default:"2048" yaml:"memory"`
	Cpus     int    `default:"1" yaml:"cpus"`
	DiskSize string `default:"10G" yaml:"disk-size"`
}

// Implement interface for setting dynamic defaults
func (c *Config) SetDefaults() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if defaults.CanUpdate(c.DataDirectory) {
		c.DataDirectory = filepath.Join(homedir, ".defnodo")
	}
	if defaults.CanUpdate(c.VolumeMounts) {
		c.VolumeMounts = []string{homedir}
	}
}

// Load a configuration from a given location.  Any unset values will
// use the default values associated.  If location does not exist, only
// default values will be used
func LoadConfig(location string) (config *Config, err error) {
	config = &Config{}
	if _, err := os.Stat(location); !os.IsNotExist(err) {
		log.Printf("Loading config from %s\n", location)
		config, err = loadConfig(config, location)
		if err != nil {
			return nil, err
		}
	} else {
		log.Printf("No config found, using all defaults...\n")
	}

	defaults.Set(config)

	configDir, err := filepath.Abs(filepath.Dir(location))
	if err != nil {
		return
	}
	config.ConfigBaseDirectory = configDir
	return
}

// Load a yaml configuration from disk
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

	config = c
	return
}

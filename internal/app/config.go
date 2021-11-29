package app

// Defnodo configuration loading and handling.
// General theory:
// Config.DataDirectory, if not an absolute path, is treated as a sibling of the config file location.
// Config.ContainerRuntime values (if filenames/paths), if not an absolute path, is assumed to be in
//   the Config.DataDirectory.
// All values set and returned have absolute paths set, all relative values have been resolved

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

// Minimal os.File interface needed for config loading
type ConfigFile interface {
	Name() string
	Read([]byte) (int, error)
}

// Config is the runtime configuration
type Config struct {
	DataDirectory       string                 `yaml:"data-directory"`
	VolumeMounts        []string               `yaml:"volume-mounts"`
	ContainerRuntime    ContainerRuntimeConfig `yaml:"container-runtime"`
	VM                  VMConfig               `yaml:"vm"`
	IsService           bool                   `default:"false"`
	Interactive         bool                   `default:"false"`
	ConfigBaseDirectory string
}

type ContainerRuntimeConfig struct {
	DaemonJson   string `default:"" yaml:"docker-daemon.json"`
	Runtime      string `default:"docker" yaml:"runtime"`
	Version      string `default:"latest" yaml:"version"`
	VersionsFile string `default:"docker.versions" yaml:"versions-file"`
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

// Load a configuration from the given file.  Any unset values will
// use the default values associated.  If no file is supplied (file == nil), only
// default values will be used
func LoadConfig(file ConfigFile) (config *Config, err error) {
	config = &Config{}

	baseDirectory, err := os.UserHomeDir()
	if err != nil {
		return
	}
	if file != nil {
		baseDirectory, err = filepath.Abs(filepath.Dir(file.Name()))
		if err != nil {
			return
		}
		err = config.load(file)
		if err != nil {
			return
		}
	}

	defaults.Set(config)
	config.ConfigBaseDirectory = baseDirectory

	config.DataDirectory = addPathPrefix(config.DataDirectory, config.ConfigBaseDirectory)
	for index, value := range config.VolumeMounts {
		config.VolumeMounts[index] = addPathPrefix(value, config.ConfigBaseDirectory)
	}
	config.ContainerRuntime.DaemonJson = addPathPrefix(config.ContainerRuntime.DaemonJson, config.DataDirectory)
	config.ContainerRuntime.VersionsFile = addPathPrefix(config.ContainerRuntime.VersionsFile, config.DataDirectory)
	return
}

// Load a yaml config file into the config
func (config *Config) load(file ConfigFile) (err error) {
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return err
	}
	return
}

// Add a prefix to the value if it is a relative rather than absolute value
func addPathPrefix(value string, prefix string) (result string) {
	result = value
	if !strings.HasPrefix(value, "/") && value != "" {
		result = filepath.Join(prefix, value)
	}
	return
}

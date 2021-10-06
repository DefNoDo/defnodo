package app

import "fmt"

func Start(config *Config) (err error) {
	fmt.Printf("Starting service with config: %+v\n", config)
	return
}

func Stop(config *Config) (err error) {
	fmt.Printf("Stopping service with config %+v\n", config)
	return
}

func Restart(config *Config) (err error) {
	err = Stop(config)
	if err != nil {
		return
	}
	err = Start(config)
	if err != nil {
		return
	}
	return
}

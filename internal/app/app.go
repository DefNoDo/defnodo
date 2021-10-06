package app

import "fmt"

type DefNoDo struct {
	config *Config
}

func NewDefNoDoService(config *Config) (defnodo *DefNoDo) {
	defnodo = &DefNoDo{
		config: config,
	}
	return
}

// Start starts the service if not already started based
// on the supplied runtime config.
func (defnodo *DefNoDo) Start() (err error) {
	fmt.Printf("Starting service with config: %+v\n", defnodo.config)
	return
}

// Stop stops the service if running based on the supplied
// runtime config.
func (defnodo *DefNoDo) Stop() (err error) {
	fmt.Printf("Stopping service with config %+v\n", defnodo.config)
	return
}

// Restart stops the service if running, then starts the service
// as defined by the supplied runtime config.
func (defnodo *DefNoDo) Restart() (err error) {
	err = defnodo.Stop()
	if err != nil {
		return
	}
	err = defnodo.Start()

	return
}

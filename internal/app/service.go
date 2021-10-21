package app

import (
	"fmt"

	"github.com/kardianos/service"
)

func ProcessService(defnodo *DefNoDo, serviceConfig *service.Config, command string) (err error) {
	// install, start, stop, restart, uninstall, special case status
	fmt.Println("Processing DefNoDo service...")
	s, err := service.New(defnodo, serviceConfig)

	if err != nil {
		return
	}

	// Special case getting the status.  It's not handled by service.Control, but should be exposed
	// to the command
	if command == "status" {
		status, err := s.Status()
		if err != nil {
			return err
		}
		fmt.Printf("Status: %+v\n", status)
		return nil
	}
	err = service.Control(s, command)

	return
}

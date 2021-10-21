package app

import (
	"fmt"
	"log"

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
		switch status {
		case service.StatusRunning:
			log.Printf("DefNoDo is running\n")
		case service.StatusStopped:
			log.Printf("DefNoDo is stopped")
		default:
			log.Printf("Unknown serice status\n")
		}
		return nil
	}
	err = service.Control(s, command)

	if err != nil {
		return
	}

	return
}

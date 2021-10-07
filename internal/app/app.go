package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

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
	linuxkitPath, err := exec.LookPath("linuxkit")
	// vpnkitPath, err = exec.LookPath("vpnkit")
	if err != nil {
		log.Print("Unable to find linuxkit binary")
		return
	}
	// ./linuxkit run hyperkit -fw bunk_uefi.fd  -disk size=4G -networking=vpnkit -vsock-ports 2376 -kernel  -data-file metadata.json -mem 4096 docker-for-mac
	err = os.Chdir(defnodo.config.DataDirectory)
	if err != nil {
		os.MkdirAll(defnodo.config.DataDirectory, 666)
	}
	cmd := exec.Command(linuxkitPath,
		"run", "hyperkit",
		"-fw", "bunk_uefi.fd",
		"-disk", "size=4G",
		"-networking", "vpnkit",
		"-vsock-ports", "2376",
		"-kernel",
		"-data-file", "metadata.json",
		"-mem", "4096",
		"-state", "defnodo-state",
		"defnodo")
	cmd.Env = os.Environ()

	log.Printf("Starting linuxkit: %v", cmd.Args)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

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

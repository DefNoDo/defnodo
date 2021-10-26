package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/ismarc/defnodo/internal/serve9p"
	"github.com/kardianos/service"
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

// Service definition interface
func (defnodo *DefNoDo) Start(s service.Service) error {
	// Start should not block, just run stuff
	go defnodo.Run()
	return nil
}

func (defnodo *DefNoDo) Stop(s service.Service) error {
	// Shutdown/terminate the service
	return nil
}

// Run the defnodo service
func (defnodo *DefNoDo) Run() (err error) {
	fmt.Printf("Starting service with config: %+v\n", defnodo.config)

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	fileShare, err := serve9p.NewServe9P("tcp://127.0.0.1:7777", []string{homedir}, true)
	if err != nil {
		log.Fatal(err)
	}
	go fileShare.Run()

	// linuxkitPath, err := exec.LookPath("linuxkit")
	// vpnkitPath, err = exec.LookPath("vpnkit")
	// if err != nil {
	// 	log.Print("Unable to find linuxkit binary")
	// 	return
	// }

	// ./linuxkit run hyperkit -fw bunk_uefi.fd  -disk size=4G -networking=vpnkit -vsock-ports 2376 -kernel  -data-file metadata.json -mem 4096 docker-for-mac
	// err = os.Chdir(defnodo.config.DataDirectory)
	// if err != nil {
	// 	os.MkdirAll(defnodo.config.DataDirectory, 666)
	// }
	cmd := exec.Command("./scripts/run_vm.sh")
	// cmd := exec.Command(linuxkitPath,
	// 	"run", "hyperkit",
	// 	"-fw", "bunk_uefi.fd",
	// 	"-disk", "size=4G",
	// 	"-networking", "vpnkit",
	// 	"-vsock-ports", "2376",
	// 	"-kernel",
	// 	"-data-file", "metadata.json",
	// 	"-mem", "4096",
	// 	"-state", "defnodo-state",
	// 	"defnodo")
	cmd.Env = os.Environ()

	log.Printf("Starting linuxkit: %v", cmd.Args)

	if !defnodo.config.IsService {
		cmd.Stdin = os.Stdin
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	return
}

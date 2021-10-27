package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

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
	fileShare, err := serve9p.NewServe9P("tcp://127.0.0.1:7777", []string{filepath.Join(homedir, "dev/defnodo")}, false)
	if err != nil {
		log.Fatal(err)
	}
	go fileShare.Run()

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	log.Printf("Path to executable: %s\n", exPath)
	linuxkitPath := filepath.Join(exPath, "linuxkit")
	log.Printf("linuxkit path: %s\n", linuxkitPath)

	// See scripts/run_vm.sh for example run command
	cmd := exec.Command(linuxkitPath,
		"run", "hyperkit",
		"-fw", filepath.Join(exPath, "..", "bunk_uefi.fd"),
		"-hyperkit", filepath.Join(exPath, "hyperkit"),
		"-vpnkit", filepath.Join(exPath, "vpnkit"),
		"-cpus", "4",
		"-mem", "8192",
		"-disk", "size=10G",
		"-networking=vpnkit",
		"-vsock-ports", "2376",
		"-squashfs",
		"-data-file", filepath.Join(exPath, "..", "metadata.json"),
		"defnodo-data/defnodo")

	cmd.Env = os.Environ()

	log.Printf("Starting linuxkit: %v", cmd.Args)

	if !defnodo.config.IsService {
		cmd.Stdin = os.Stdin
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	defer os.Remove(filepath.Join(exPath, "defnodo-data/defnodo-state/hyperkit.pid"))
	cmd.Run()

	return
}

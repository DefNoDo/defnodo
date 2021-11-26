package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

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
	log.Printf("Starting service with config: %+v\n", defnodo.config)

	fileShare, err := serve9p.NewServe9P("tcp://127.0.0.1:7777", defnodo.config.VolumeMounts, false)
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

	// Create data directory if it doesn't exist
	err = os.MkdirAll(defnodo.config.DataDirectory, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	runLocation := defnodo.config.DataDirectory

	dataFile, err := defnodo.generateMetadata(runLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dataFile)

	// Running vpnkit separately stuff
	// ethSock := filepath.Join(runLocation, "defnodo-state", "vpnkit_eth.sock")
	// portSock := filepath.Join(runLocation, "defnodo-state", "vpnkit_port.sock")
	// vsockSock := filepath.Join(runLocation, "defnodo-state", "connect")
	// go defnodo.runVPNKit(filepath.Join(exPath, "vpnkit"),
	// 	ethSock,
	// 	portSock,
	// 	vsockSock)

	// networkingParam := fmt.Sprintf("-networking=vpnkit,%s,%s", ethSock, portSock)

	// See scripts/run_vm.sh for example run command
	cmd := exec.Command(linuxkitPath,
		"run", "hyperkit",
		"-fw", filepath.Join(exPath, "..", "bunk_uefi.fd"),
		"-hyperkit", filepath.Join(exPath, "hyperkit"),
		"-vpnkit", filepath.Join(exPath, "vpnkit"),
		"-cpus", strconv.Itoa(defnodo.config.VM.Cpus),
		"-mem", strconv.Itoa(defnodo.config.VM.Memory),
		"-disk", fmt.Sprintf("size=%s", defnodo.config.VM.DiskSize),
		"-networking=vpnkit",
		// networkingParam,
		"-vsock-ports", "2376",
		"-squashfs",
		"-data-file", dataFile,
		filepath.Join(runLocation, "defnodo"))

	cmd.Env = os.Environ()

	log.Printf("Starting linuxkit: %v", cmd.Args)

	if defnodo.config.Interactive {
		cmd.Stdin = os.Stdin
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	startCleanupHandler(exPath, dataFile)
	cmd.Run()

	return
}

func (defnodo *DefNoDo) runVPNKit(vpnkitBin string, ethSock string, portSock string, vsockSock string) {
	cmd := exec.Command(vpnkitBin,
		"--ethernet", ethSock,
		"--port", portSock,
		"--vsock-path", vsockSock)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = os.Environ()

	log.Printf("Staring vpnkit: %+v\n", cmd.Args)
	cmd.Run()
}

func startCleanupHandler(basePath string, tempFile string) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Printf("Removing file: %s\n", filepath.Join(basePath, "..", "defnodo-data/defnodo-state/hyperkit.pid"))
		os.Remove(filepath.Join(basePath, "..", "defnodo-data/defnodo-state/hyperkit.pid"))
		os.Remove(tempFile)
		os.Exit(0)
	}()
}

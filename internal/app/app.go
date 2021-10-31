package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
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

	fileShare, err := serve9p.NewServe9P("tcp://127.0.0.1:7777", defnodo.config.VolumeMounts, defnodo.config.Interactive)
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

	runLocation := defnodo.config.DataDirectory
	if !strings.HasPrefix(runLocation, "/") {
		runLocation = filepath.Join(exPath, "..", defnodo.config.DataDirectory)
	}

	dataFile, err := defnodo.generateMetadata(runLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dataFile)

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

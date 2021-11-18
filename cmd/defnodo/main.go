package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ismarc/defnodo/internal/app"
	"github.com/kardianos/service"
	"github.com/urfave/cli/v2"
)

type GlobalConfig struct {
	Config        *app.Config
	ServiceConfig *service.Config
	Interactive   bool
}

func main() {
	globalFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Load defnodo configuration from `FILE`",
			Aliases: []string{"c"},
			Value:   fmt.Sprintf("%s/.defnodorc", getUserDir()),
			EnvVars: []string{"DEFNODORC"},
		},
	}
	runFlags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "interactive",
			Usage:   "Run defnodo with a local interactive shell",
			Aliases: []string{"i"},
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "service",
			Usage:   "Run defnodo in a way suitable for a service",
			Value:   false,
			EnvVars: []string{"DEFNODO_SERVICE"},
		},
	}
	serviceFlags := []cli.Flag{}
	serviceInstallFlags := []cli.Flag{}
	serviceStartFlags := []cli.Flag{}
	serviceStopFlags := []cli.Flag{}
	serviceRestartFlags := []cli.Flag{}
	serviceUninstallFlags := []cli.Flag{}
	serviceStatusFlags := []cli.Flag{}

	commands := []*cli.Command{
		{
			Name:    "run",
			Usage:   "Run defnodo and the underlying docker service",
			Aliases: []string{"r"},
			Flags:   runFlags,
			Action: func(c *cli.Context) (err error) {
				globalConfig, err := loadGlobalOptions(c)
				if err != nil {
					log.Fatal(err)
				}
				defnodo := app.NewDefNoDoService(globalConfig.Config)
				err = defnodo.Run()
				if err != nil {
					log.Fatal(err)
				}
				return
			},
		},
		{
			Name:    "service",
			Usage:   "Control defnodo service",
			Aliases: []string{"s"},
			Flags:   serviceFlags,
			Subcommands: []*cli.Command{
				{
					Name:    "status",
					Usage:   "Display defnodo service status",
					Aliases: []string{"t"},
					Flags:   serviceStatusFlags,
					Action: func(c *cli.Context) (err error) {
						globalConfig, err := loadGlobalOptions(c)
						if err != nil {
							log.Fatal(err)
						}
						defnodo := app.NewDefNoDoService(globalConfig.Config)
						err = app.ProcessService(defnodo, globalConfig.ServiceConfig, "status")
						if err != nil {
							log.Fatal(err)
						}
						return
					},
				},
				{
					Name:    "install",
					Usage:   "Install the defnodo service",
					Aliases: []string{"i"},
					Flags:   serviceInstallFlags,
					Action: func(c *cli.Context) (err error) {
						globalConfig, err := loadGlobalOptions(c)
						if err != nil {
							log.Fatal(err)
						}
						defnodo := app.NewDefNoDoService(globalConfig.Config)
						err = app.ProcessService(defnodo, globalConfig.ServiceConfig, "install")
						if err != nil {
							log.Fatal(err)
						}
						return
					},
				},
				{
					Name:    "uninstall",
					Usage:   "Uninstall the defnodo service",
					Aliases: []string{"u"},
					Flags:   serviceUninstallFlags,
					Action: func(c *cli.Context) (err error) {
						globalConfig, err := loadGlobalOptions(c)
						if err != nil {
							log.Fatal(err)
						}
						defnodo := app.NewDefNoDoService(globalConfig.Config)
						err = app.ProcessService(defnodo, globalConfig.ServiceConfig, "uninstall")
						if err != nil {
							log.Fatal(err)
						}
						return
					},
				},
				{
					Name:    "start",
					Usage:   "Start the defnodo service",
					Aliases: []string{"s"},
					Flags:   serviceStartFlags,
					Action: func(c *cli.Context) (err error) {
						globalConfig, err := loadGlobalOptions(c)
						if err != nil {
							log.Fatal(err)
						}
						defnodo := app.NewDefNoDoService(globalConfig.Config)
						err = app.ProcessService(defnodo, globalConfig.ServiceConfig, "start")
						if err != nil {
							log.Fatal(err)
						}
						return
					},
				},
				{
					Name:    "stop",
					Usage:   "Stop the defnodo service",
					Aliases: []string{"h"},
					Flags:   serviceStopFlags,
					Action: func(c *cli.Context) (err error) {
						globalConfig, err := loadGlobalOptions(c)
						if err != nil {
							log.Fatal(err)
						}
						defnodo := app.NewDefNoDoService(globalConfig.Config)
						err = app.ProcessService(defnodo, globalConfig.ServiceConfig, "stop")
						if err != nil {
							log.Fatal(err)
						}
						return
					},
				},
				{
					Name:    "restart",
					Usage:   "Restart the defnodo service",
					Aliases: []string{"r"},
					Flags:   serviceRestartFlags,
					Action: func(c *cli.Context) (err error) {
						globalConfig, err := loadGlobalOptions(c)
						if err != nil {
							log.Fatal(err)
						}
						defnodo := app.NewDefNoDoService(globalConfig.Config)
						err = app.ProcessService(defnodo, globalConfig.ServiceConfig, "restart")
						if err != nil {
							log.Fatal(err)
						}
						return
					},
				},
			},
		},
	}
	app := &cli.App{
		Name:                 "DefNoDo",
		HelpName:             "defnodo",
		Commands:             commands,
		Flags:                globalFlags,
		EnableBashCompletion: true,
		HideHelp:             false,
		HideHelpCommand:      false,
		HideVersion:          false,
		Compiled:             time.Time{},
		Version:              "v0.1",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadGlobalOptions(c *cli.Context) (result *GlobalConfig, err error) {
	config, err := app.LoadConfig(c.String("config"))
	if err != nil {
		return
	}
	config.IsService = c.Bool("service")
	config.Interactive = c.Bool("interactive")

	serviceOptions := make(service.KeyValue)
	serviceOptions["UserService"] = true
	configLocation, err := filepath.Abs(c.String("config"))
	if err != nil {
		return
	}
	serviceConfig := &service.Config{
		Name:        "defnodo",
		DisplayName: "DefNoDo",
		Description: "DefNoDO service for running dockerd and podman on MacOS",
		Option:      serviceOptions,
		Arguments: []string{
			"-c",
			configLocation,
			"run",
			"--service",
		},
	}
	result = &GlobalConfig{
		Config:        config,
		ServiceConfig: serviceConfig,
		Interactive:   c.Bool("interactive"),
	}
	return
}

func getUserDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "/"
	}
	return dirname
}

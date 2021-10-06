package main

import (
	"log"
	"os"
	"time"

	"github.com/ismarc/defnodo/internal/app"
	"github.com/urfave/cli/v2"
)

type GlobalConfig struct {
	Config *app.Config
}

func main() {
	globalFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Load defnodo configuration from `FILE`",
			Aliases: []string{"c"},
			Value:   "~/.defnodorc",
			EnvVars: []string{"DEFNODORC"},
		},
	}
	startFlags := []cli.Flag{}
	stopFlags := []cli.Flag{}
	restartFlags := []cli.Flag{}

	commands := []*cli.Command{
		{
			Name:    "start",
			Usage:   "Start the defnodo and underlying docker service",
			Aliases: []string{"s"},
			Flags:   startFlags,
			Action: func(c *cli.Context) (err error) {
				globalConfig, err := loadGlobalOptions(c)
				if err != nil {
					log.Fatal(err)
				}
				err = app.Start(globalConfig.Config)
				if err != nil {
					log.Fatal(err)
				}
				return
			},
		},
		{
			Name:    "stop",
			Usage:   "Stop the defnodo and underlying docker service",
			Aliases: []string{"t"},
			Flags:   stopFlags,
			Action: func(c *cli.Context) (err error) {
				globalConfig, err := loadGlobalOptions(c)
				if err != nil {
					log.Fatal(err)
				}

				err = app.Stop(globalConfig.Config)
				if err != nil {
					log.Fatal(err)
				}
				return
			},
		},
		{
			Name:    "restart",
			Usage:   "Restart the defnodo and underlying docker service",
			Aliases: []string{"r"},
			Flags:   restartFlags,
			Action: func(c *cli.Context) (err error) {
				globalConfig, err := loadGlobalOptions(c)
				if err != nil {
					log.Fatal(err)
				}

				err = app.Restart(globalConfig.Config)
				if err != nil {
					log.Fatal(err)
				}
				return
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
	result = &GlobalConfig{
		Config: config,
	}
	return
}

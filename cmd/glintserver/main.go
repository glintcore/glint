package main

import (
	"os"

	"github.com/nassibnassar/goconfig/ini"
	"github.com/urfave/cli"
)

var glintconfig *ini.Config
var glintbaseurl string

func main() {
	var err error
	// Run commands specified on the command line.
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version",
		Usage: "print the Glint version",
	}
	var app = cli.NewApp()
	app.Name = "glintserver"
	app.Version = "0"
	app.HideVersion = true
	app.HelpName = "glintserver"
	app.Usage = "Glint server for sharing and integrating data sets"
	app.UsageText = "glintserver [command] [arguments]"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		/*
			// Verbose flag not currently implemented.
			cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "enable verbose output",
			},
		*/
		cli.BoolFlag{
			Name:   "debug",
			Hidden: true,
			Usage:  "enable debugging output to log",
		},
		cli.StringFlag{
			Name:  "config-file",
			Usage: "server configuration file",
		},
	}
	app.Commands = []cli.Command{
		/*
			cli.Command{
				Name:      "run-old",
				Hidden:    true,
				Usage:     "Runs the old server",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "port",
						Usage: "port number to listen on",
					},
					cli.BoolFlag{
						Name:   "debug-allow-insecure-cors",
						Hidden: true,
					},
				},
				Action: func(c *cli.Context) error {
					err = cliOldServer(c)
					if err != nil {
						return cli.NewExitError("Exited with error", 1)
					}
					return nil
				},
			},
		*/
		cli.Command{
			Name:      "run",
			Usage:     "Runs the server",
			ArgsUsage: " ",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port",
					Usage: "port number to listen on",
				},
				cli.BoolFlag{
					Name:   "debug-allow-insecure-cors",
					Hidden: true,
				},
			},
			Action: func(c *cli.Context) error {
				if err := cliNewServer(c); err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
		cli.Command{
			Name:      "adduser",
			Usage:     "Adds a new user",
			ArgsUsage: " ",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "user",
					Usage: "new username",
				},
				cli.StringFlag{
					Name:  "fullname",
					Usage: "full name of user",
				},
				cli.StringFlag{
					Name:  "email",
					Usage: "email address of user",
				},
			},
			Action: func(c *cli.Context) error {
				err = cliAddUser(c)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
		cli.Command{
			Name:      "passwd",
			Usage:     "Changes a user's password",
			ArgsUsage: " ",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "user",
					Usage: "username",
				},
			},
			Action: func(c *cli.Context) error {
				err = cliPasswd(c)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}

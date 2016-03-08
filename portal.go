package main

import (
	"os"

	"github.com/codegangsta/cli"

	"github.com/njpatel/portal/receive"
	"github.com/njpatel/portal/send"
	"github.com/njpatel/portal/server"
)

var (
	address  string
	author   = cli.Author{Name: "Neil Jagdish Patel", Email: "njpatel@gmail.com"}
	force    bool
	insecure bool
	secret   string
	version  string
)

func main() {
	app := cli.NewApp()
	app.Name = "portal"
	app.Usage = "Easily send one or more files from A to B. https://github.com/njpatel/portal"
	app.Version = version
	app.Authors = []cli.Author{author}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "address, a",
			Value:       "localhost:3421",
			Usage:       "portal server address",
			Destination: &address,
			EnvVar:      "PORTAL_HOST",
		},

		cli.BoolFlag{
			Name:        "insecure, i",
			Usage:       "connect to portal server insecurely",
			Destination: &insecure,
			EnvVar:      "PORTAL_INSECURE",
		},

		cli.StringFlag{
			Name:        "secret",
			Value:       "",
			Usage:       "portal server shared secret",
			Destination: &secret,
			EnvVar:      "PORTAL_SECRET",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Start the portal server",
			Action: func(c *cli.Context) {
				server.Run(&server.Config{
					Address: address,
					Secret:  secret,
				})
			},
		},
		{
			Name:      "send",
			Usage:     "Send file(s) to another machine",
			ArgsUsage: "<file1> [dir2] [file3]",
			Action: func(c *cli.Context) {
				send.Run(&send.Config{
					Address:  address,
					Insecure: insecure,
					Secret:   secret,
				}, c.Args())
			},
		},
		{
			Name:      "sync",
			Usage:     "Send file(s) to another machine, and keep them in sync",
			ArgsUsage: "<file1> [dir2] [file3]",
			Action: func(c *cli.Context) {
				send.RunSync(&send.Config{
					Address:  address,
					Insecure: insecure,
					Secret:   secret,
				}, c.Args())
			},
		},
		{
			Name:      "receive",
			Aliases:   []string{"get"},
			Usage:     "Receive files from another machine",
			ArgsUsage: "<token> <output dir=./>",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "force, f",
					Usage:       "overwrite files",
					Destination: &force,
					EnvVar:      "PORTAL_FORCE",
				},
			},
			Action: func(c *cli.Context) {
				receive.Run(&receive.Config{
					Address:  address,
					Insecure: insecure,
					Secret:   secret,
					Force:    force,
				}, c.Args().Get(0), c.Args().Get(1))
			},
		},
	}

	app.Run(os.Args)
}

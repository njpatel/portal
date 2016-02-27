package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var (
	author   = cli.Author{Name: "Neil Jagdish Patel", Email: "njpatel@gmail.com"}
	server   string
	insecure bool
	secret   string
	version  string
)

func main() {
	app := cli.NewApp()
	app.Name = "portal"
	app.Version = version
	app.Authors = []cli.Author{author}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "server, s",
			Value:       "portal.njp.io:3421",
			Usage:       "portal server hostname",
			Destination: &server,
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
			Usage: "start the portal server",
			Action: func(c *cli.Context) {
				fmt.Println("Serve!")
			},
		},
		{
			Name:  "send",
			Usage: "send <file1> [file2] [file3]",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					fmt.Println("Need at least one file to send")
					return
				}
				for _, f := range c.Args() {
					fmt.Println("Sending " + f)
				}
			},
		},
		{
			Name:  "sync",
			Usage: "sync <file1> [file2] [file3]",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					fmt.Println("Need at least one file to sync")
					return
				}
				for _, f := range c.Args() {
					fmt.Println("Sending " + f)
				}
			},
		},
		{
			Name:    "receive",
			Aliases: []string{"get"},
			Usage:   "receive <token> <output dir=./>",
			Action: func(c *cli.Context) {
				if len(c.Args()) < 1 {
					fmt.Println("Need the <token> and <output dir>")
					return
				}
				fmt.Printf("token:%s output:%s\n", c.Args().First(), c.Args().Get(1))
			},
		},
	}

	app.Run(os.Args)
}

package main

import (
	"fmt"
	"github.com/ModelRocket/coprhd"
	"github.com/codegangsta/cli"
	"os"
	"runtime"
)

const (
	BuildVersion = "0.1"
)

func main() {
	app := cli.NewApp()

	app.Version = BuildVersion
	app.Copyright = "ModelRocket 2015, 2016"
	app.Authors = []cli.Author{{"Rob Rodriguez", "rob@rocketlabs.us"}}
	app.Name = "corptop"
	app.Usage = "coprhd cli tool"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "user, u",
			Usage: "set username",
		},
		cli.StringFlag{
			Name:  "pass, p",
			Usage: "set password",
		},
		cli.StringFlag{
			Name:  "host, H",
			Usage: "set server address",
		},
		cli.IntFlag{
			Name:  "port, P",
			Value: 4443,
			Usage: "set server port",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Display version",
			Action:  showVersion,
		},
		{
			Name:    "token",
			Aliases: []string{"t"},
			Usage:   "Get user proxy authentication token",
			Action:  getToken,
		},
	}

	app.Action = cli.ShowAppHelp

	app.Run(os.Args)
}

func showVersion(c *cli.Context) {
	fmt.Println("coprtop version:", BuildVersion)
	fmt.Println("go version:", runtime.Version())
	fmt.Println("OS:", runtime.GOOS)
	fmt.Println("arch:", runtime.GOARCH)
}

func getToken(c *cli.Context) {
	user := c.GlobalString("user")
	pass := c.GlobalString("pass")
	host := c.GlobalString("host")
	port := c.GlobalString("port")

	path := fmt.Sprintf("https://%s:%s/", host, port)

	fmt.Printf("\nGetting proxy token for %s => %s...\n\n", user, path)

	token, err := coprhd.GetProxyToken(path, user, pass)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
	fmt.Printf("%s\n\n", token)
}

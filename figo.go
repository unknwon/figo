package main

import (
	"os"

	"github.com/codegangsta/cli"

	"github.com/Unknwon/figo/cmd"
)

const APP_VER = "0.0.0.0929"

func main() {
	app := cli.NewApp()
	app.Name = "Figo"
	app.Usage = "Fig in Go"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.CmdBuild,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}

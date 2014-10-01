package main

import (
	"os"

	"github.com/codegangsta/cli"

	"github.com/Unknwon/figo/cmd"
	"github.com/Unknwon/figo/modules/base"
)

const APP_VER = "0.0.0.0930"

func init() {
	base.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "Figo"
	app.Usage = "Fig in Go"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.CmdBuild,
	}
	app.Flags = append(app.Flags, []cli.Flag{
		cli.StringFlag{"file, f", "fig.yml", "Specify an alternate fig file (default: fig.yml)", "FIG_FILE"},
		cli.StringFlag{"project-name, p", "", "Specify an alternate project name (default: directory name)", ""},
		cli.BoolFlag{"verbose", "Show more output", ""},
	}...)
	app.Run(os.Args)
}

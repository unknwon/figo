package main

import (
	"os"

	"github.com/codegangsta/cli"

	"github.com/Unknwon/figo/cmd"
	"github.com/Unknwon/figo/modules/base"
)

const APP_VER = "0.2.0.1020"

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
		cmd.CmdKill,
		//cmd.CmdLogs,
		cmd.CmdPort,
		//cmd.CmdPs,
		cmd.CmdPull,
		//cmd.Rm,
		//cmd.Run,
		//cmd.Scale,
		cmd.CmdStart,
		cmd.CmdStop,
		cmd.CmdRestart,
		cmd.CmdUp,
	}
	app.Flags = append(app.Flags, []cli.Flag{
		cli.StringFlag{"file, f", "fig.yml", "Specify an alternate fig file (default: fig.yml)", "FIG_FILE"},
		cli.StringFlag{"project-name, p", "", "Specify an alternate project name (default: directory name)", ""},
	}...)
	app.Run(os.Args)
}

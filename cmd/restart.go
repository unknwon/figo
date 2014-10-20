package cmd

import (
	"github.com/codegangsta/cli"

	"github.com/Unknwon/figo/modules/log"
)

var CmdRestart = cli.Command{
	Name:        "restart",
	Usage:       "Restart running containers",
	Description: `Restart running containers.`,
	Action:      runRestart,
	Flags: []cli.Flag{
		cli.BoolFlag{"verbose, v", "Show more output", ""},
	},
}

func runRestart(ctx *cli.Context) {
	pro, err := setup(ctx)
	if err != nil {
		log.Fatal("%v", err)
	}

	if err := pro.Restart(ctx.Args()); err != nil {
		log.Fatal("Fail to restart project: %v", err)
	}
}

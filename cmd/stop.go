package cmd

import (
	"github.com/codegangsta/cli"

	"github.com/Unknwon/figo/modules/log"
)

var CmdStop = cli.Command{
	Name:  "stop",
	Usage: "Stop running containers without removing them",
	Description: `Stop running containers without removing them.
They can be started again with 'fig start'.`,
	Action: runStop,
	Flags: []cli.Flag{
		cli.BoolFlag{"verbose, v", "Show more output", ""},
	},
}

func runStop(ctx *cli.Context) {
	pro, err := setup(ctx)
	if err != nil {
		log.Fatal("%v", err)
	}

	if err := pro.Stop(ctx.Args()); err != nil {
		log.Fatal("Fail to stop project: %v", err)
	}
}

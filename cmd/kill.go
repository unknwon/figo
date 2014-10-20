package cmd

import (
	"github.com/codegangsta/cli"

	"github.com/Unknwon/figo/modules/log"
)

var CmdKill = cli.Command{
	Name:        "kill",
	Usage:       "Force stop service containers",
	Description: `Force stop service containers.`,
	Action:      runKill,
	Flags: []cli.Flag{
		cli.BoolFlag{"verbose, v", "Show more output", ""},
	},
}

func runKill(ctx *cli.Context) {
	pro, err := setup(ctx)
	if err != nil {
		log.Fatal("%v", err)
	}

	if err := pro.Kill(ctx.Args()); err != nil {
		log.Fatal("Fail to kill project: %v", err)
	}
}

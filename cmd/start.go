package cmd

import (
	"github.com/codegangsta/cli"

	"github.com/Unknwon/figo/modules/log"
)

var CmdStart = cli.Command{
	Name:        "start",
	Usage:       "Start existing containers for a service",
	Description: `Start existing containers for a service.`,
	Action:      runStart,
	Flags: []cli.Flag{
		cli.BoolFlag{"verbose, v", "Show more output", ""},
	},
}

func runStart(ctx *cli.Context) {
	pro, err := setup(ctx)
	if err != nil {
		log.Fatal("%v", err)
	}

	if err := pro.Start(ctx.Args()); err != nil {
		log.Fatal("Fail to start project: %v", err)
	}
}

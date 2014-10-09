package cmd

import (
	"github.com/codegangsta/cli"

	"github.com/Unknwon/figo/modules/log"
)

var CmdStart = cli.Command{
	Name:        "start",
	Usage:       "Start existing containers.",
	Description: `Start existing containers.`,
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

	if err := pro.Build(ctx.Args(), ctx.Bool("no-cache")); err != nil {
		log.Fatal("Fail to build project: %v", err)
	}
}

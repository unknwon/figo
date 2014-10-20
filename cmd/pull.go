package cmd

import (
	"github.com/codegangsta/cli"

	"github.com/Unknwon/figo/modules/log"
)

var CmdPull = cli.Command{
	Name:        "pull",
	Usage:       "Pulls images for services",
	Description: `Pulls images for services.`,
	Action:      runPull,
	Flags: []cli.Flag{
		cli.BoolFlag{"allow-insecure-ssl", "Allow insecure connections to the docker registry", ""},
		cli.BoolFlag{"verbose, v", "Show more output", ""},
	},
}

func runPull(ctx *cli.Context) {
	pro, err := setup(ctx)
	if err != nil {
		log.Fatal("%v", err)
	}

	if err = pro.Pull(ctx.Args(), ctx.Bool("allow-insecure-ssl")); err != nil {
		log.Fatal("Fail to pull projetc: %v", err)
	}
}

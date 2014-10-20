package cmd

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/Unknwon/figo/modules/log"
)

var CmdPort = cli.Command{
	Name:        "port",
	Usage:       "Print the public port for a port binding",
	Description: `Print the public port for a port binding.`,
	Action:      runPort,
	Flags: []cli.Flag{
		cli.StringFlag{"protocol", "tcp", "tcp or udp (defaults to tcp)", ""},
		cli.IntFlag{"index", 1, "index of the container if there are multiple instances of a service (defaults to 1)", ""},
		cli.IntFlag{"private-port", 0, "", ""},
		cli.BoolFlag{"verbose, v", "Show more output", ""},
	},
}

func runPort(ctx *cli.Context) {
	pro, err := setup(ctx)
	if err != nil {
		log.Fatal("%v", err)
	}

	service, err := pro.GetService(ctx.Args().First())
	if err != nil {
		log.Fatal("Fail to get service(%s): %v", ctx.Args().First(), err)
	}
	container, err := service.GetContainer(ctx.Int("index"))
	if err != nil {
		log.Fatal("Fail to get container(%d): %v", ctx.Int("index"), err)
	}

	fmt.Println(container.GetLocalPort(ctx.Int("private-port"), ctx.String("protocol")))
}

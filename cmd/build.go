package cmd

import (
	"github.com/codegangsta/cli"

	"github.com/Unknwon/figo/modules/log"
)

var CmdBuild = cli.Command{
	Name:  "build",
	Usage: "Build or rebuild services",
	Description: `Build or rebuild services.

Services are built once and then tagged as "project_service",
e.g. "figtest_db". If you change a service's "Dockerfile" or the
contents of its build directory, you can run "fig build" to rebuild it.`,
	Action: runBuild,
	Flags: []cli.Flag{
		cli.BoolFlag{"no-cache", "Do not use cache when building the image", ""},
	},
}

func runBuild(ctx *cli.Context) {
	pro, err := Setup(ctx)
	if err != nil {
		log.Fatal("%v", err)
	}
	// endpoint := "unix:///var/run/docker.sock"
	if err := pro.Build(ctx.Args(), ctx.Bool("no-cache")); err != nil {
		log.Fatal("Fail to build project: %v", err)
	}
}

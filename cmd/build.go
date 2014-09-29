package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"

	"github.com/Unknwon/figo/modules/core"
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
	// FIXME: init command context.
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	pro := core.NewProject("hi", []*core.Service{}, client)
	if err := pro.Build([]string{}, false); err != nil {
		log.Fatal("Fail to build project: %v", err)
	}
}

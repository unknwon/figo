package cmd

import (
	"github.com/codegangsta/cli"

	"github.com/Unknwon/figo/modules/log"
)

var CmdUp = cli.Command{
	Name:  "up",
	Usage: "Create and start containers",
	Description: `Build, (re)create, start and attach to containers for a service.

By default, 'fig up' will aggregate the output of each container, and
when it exits, all containers will be stopped. If you run 'fig up -d',
it'll start the containers in the background and leave them running.

If there are existing containers for a service, 'fig up' will stop
and recreate them (preserving mounted volumes with volumes-from),
so that changes in 'fig.yml' are picked up. If you do not want existing
containers to be recreated, 'fig up --no-recreate' will re-use existing
containers.`,
	Action: runUp,
	Flags: []cli.Flag{
		cli.BoolFlag{"d", "Detached mode: Run containers in the background, print new container names.", ""},
		cli.BoolFlag{"no-deps", "Don't start linked services", ""},
		cli.BoolFlag{"no-recreate", "If containers already exist, don't recreate them", ""},
		cli.BoolFlag{"verbose, v", "Show more output", ""},
	},
}

func runUp(ctx *cli.Context) {
	pro, err := setup(ctx)
	if err != nil {
		log.Fatal("%v", err)
	}

	if err := pro.Up(ctx.Args(), !ctx.Bool("no-deps"), !ctx.Bool("no-recreate")); err != nil {
		log.Fatal("Fail to up project: %v", err)
	}

	// FIXME: catch ctrl+c
	// s in project.get_services(service_names) for c in s.containers()]

	//        if not detached:
	//            print("Attaching to", list_containers(to_attach))
	//            log_printer = LogPrinter(to_attach, attach_params={"logs": True}, monochrome=monochrome)

	//            try:
	//                log_printer.run()
	//            finally:
	//                def handler(signal, frame):
	//                    project.kill(service_names=service_names)
	//                    sys.exit(0)

	//                signal.signal(signal.SIGINT, handler)

	//                print("Gracefully stopping... (press Ctrl+C again to force)")
	//                project.stop(service_names=service_names)
}

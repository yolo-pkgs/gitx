package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

const defaultExecTimeout = 10 * time.Second

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "rapid-push",
				Usage: "wait for commit and push as soon as one is available",
				Action: func(cCtx *cli.Context) error {
					return rapidPush()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Panic(err)
	}
}

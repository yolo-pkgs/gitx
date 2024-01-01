package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/yolo-pkgs/grace"
)

const defaultExecTimeout = 10 * time.Second

func notifySend(msg string) {
	_, _ = grace.RunShTimed(fmt.Sprintf(`notify-send '%s'`, msg), defaultExecTimeout)
}

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
		notifySend(err.Error())
		log.Panic(err)
	}
}

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
	_, _ = grace.RunTimedSh(defaultExecTimeout, fmt.Sprintf("notify-send '%s'", msg))
}

func main() {
	app := &cli.App{
		Usage: `Wildly unstable functions for git`,
		Commands: []*cli.Command{
			{
				Name:  "rapid-push",
				Usage: "wait for commit and push as soon as one is available",
				Action: func(cCtx *cli.Context) error {
					return rapidPush()
				},
			},
			{
				Name:    "branch",
				Usage:   "create branch",
				Aliases: []string{"b"},
				Action: func(c *cli.Context) error {
					fromDefault, err := fromDefaultBranch()
					if err != nil {
						return err
					}

					gid, err := writeNewBranchGID()
					if err != nil {
						return err
					}

					branchName := c.Args().First()
					if branchName == "" {
						return createRandomBranch(gid, fromDefault, false)
					}

					return createGlobalBranch(gid, branchName, fromDefault, false)
				},
			},
			{
				Name:  "randx",
				Usage: "create random experimental branch",
				Action: func(_ *cli.Context) error {
					fromDefault, err := fromDefaultBranch()
					if err != nil {
						return err
					}

					gid, err := writeNewBranchGID()
					if err != nil {
						return err
					}
					return createRandomBranch(gid, fromDefault, true)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		notifySend("error executing gitx")
		log.Panic(err)
	}
}

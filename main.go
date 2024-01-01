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
				Usage:   "branch functions",
				Aliases: []string{"b"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "random",
						Aliases: []string{"r"},
						Usage:   "create a branch with generated random name",
					},
				},
				Action: func(c *cli.Context) error {
					if err := onlyFromDefaultBranch(); err != nil {
						return err
					}

					if c.Bool("random") {
						gid, err := writeNewBranchGID()
						if err != nil {
							return err
						}
						return createRandomBranch(gid)
					}
					branchName := c.Args().First()
					if branchName == "" {
						return listBranches()
					}

					gid, err := writeNewBranchGID()
					if err != nil {
						return err
					}
					return createGlobalBranch(gid, branchName)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		notifySend("error executing gitx")
		log.Panic(err)
	}
}

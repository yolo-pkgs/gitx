package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/yolo-pkgs/grace"

	"github.com/yolo-pkgs/gitx/status"
	"github.com/yolo-pkgs/gitx/tag"
)

const defaultExecTimeout = 10 * time.Second

func notifySend(msg string) {
	_, _ = grace.RunTimedSh(defaultExecTimeout, fmt.Sprintf("notify-send -a gitx '%s'", msg))
}

func main() {
	app := &cli.App{
		Usage: `Wildly unstable functions for git`,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "notify", Aliases: []string{"n"}, Usage: "notify instead of stdout"},
		},
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
						return createRandomBranch(gid, fromDefault)
					}

					return createGlobalBranch(gid, branchName, fromDefault)
				},
			},
			{
				Name:    "status",
				Usage:   "cool status",
				Aliases: []string{"s"},
				Action: func(c *cli.Context) error {
					statusMsg, err := status.CoolStatus()
					if err != nil {
						return err
					}
					if c.Bool("notify") {
						notifySend(statusMsg)
					} else {
						fmt.Println(statusMsg)
					}
					return nil
				},
			},
			{
				Name:    "tagnew",
				Usage:   "create new version tag",
				Aliases: []string{"tn"},
				Action: func(c *cli.Context) error {
					if err := tag.Patch(false); err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:    "stage",
				Usage:   "create new stage tag",
				Aliases: []string{"st"},
				Action: func(c *cli.Context) error {
					if err := tag.Patch(true); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		notifySend("error executing gitx")
		log.Panic(err)
	}
}

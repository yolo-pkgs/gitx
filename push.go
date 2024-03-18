package main

import (
	"fmt"
	"time"

	"github.com/yolo-pkgs/grace"

	"github.com/yolo-pkgs/gitx/generic"
)

const (
	waitForCommitSeconds = 5
	watchForCommitsInLastSeconds = 180
)

func rapidPush() error {
	deadline := time.NewTimer(waitForCommitSeconds * time.Second)

	for {
		select {
		case <-deadline.C:
			return nil
		default:
		}

		lastCommitTime, err := generic.LastCommitUnixtime()
		if err != nil {
			return err
		}

		if time.Now().Unix()-lastCommitTime < watchForCommitsInLastSeconds {
			if err := gitPush(); err != nil {
				return fmt.Errorf("failed to push: %w", err)
			}

			notifySend("pushed")

			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func gitPush() error {
	_, err := grace.RunTimed(20*time.Second, nil, "git", "push", "--quiet")
	return err
}

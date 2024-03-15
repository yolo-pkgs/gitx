package main

import (
	"fmt"
	"time"

	"github.com/yolo-pkgs/grace"

	"github.com/yolo-pkgs/gitx/generic"
)

const deadlineTimeoutSeconds = 6

func rapidPush() error {
	deadline := time.NewTimer(deadlineTimeoutSeconds * time.Second)

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

		if time.Now().Unix()-lastCommitTime < deadlineTimeoutSeconds+1 {
			if err := gitPush(); err != nil {
				return fmt.Errorf("failed to push: %w", err)
			}

			notifySend("pushed")

			return nil
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func gitPush() error {
	_, err := grace.RunTimed(20*time.Second, "git", "push", "--quiet")
	return err
}

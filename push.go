package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yolo-pkgs/grace"
)

const deadlineTimeoutSeconds = 6

func lastCommitUnixtime() (int64, error) {
	output, err := grace.RunTimed(defaultExecTimeout, "git", "log", "-1", "--format=%ct")
	if err != nil {
		return 0, err
	}

	trimmed := strings.TrimSpace(output)

	timestamp, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return 0, err
	}

	return timestamp, nil
}

func rapidPush() error {
	deadline := time.NewTimer(deadlineTimeoutSeconds * time.Second)

	for {
		select {
		case <-deadline.C:
			return nil
		default:
		}

		lastCommitTime, err := lastCommitUnixtime()
		if err != nil {
			return err
		}

		if time.Now().Unix()-lastCommitTime < deadlineTimeoutSeconds+1 {
			if err := gitPush(); err != nil {
				return fmt.Errorf("failed to push: %w", err)
			}

			notifySend("git pushed!")

			return nil
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func gitPush() error {
	_, err := grace.RunTimed(20*time.Second, "git", "push", "--quiet")
	return err
}

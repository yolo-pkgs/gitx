package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/yolo-pkgs/grace"

	"github.com/yolo-pkgs/gitx/git"
)

const (
	localGitTimeout              = 3 * time.Second
	pushTimeout                  = 90 * time.Second
	waitForCommit                = 3 * time.Second
	watchForCommitsInLastSeconds = 60
)

func rapidPush() error {
	deadline := time.NewTimer(waitForCommit)

	currentBranch, err := git.CurrentBranch(localGitTimeout)
	if err != nil {
		notifySend("rp: failed to get current branch")

		return err
	}

	for {
		select {
		case <-deadline.C:
			notifySend(fmt.Sprintf("rp: found no commits (waited %d seconds)", waitForCommit))

			return nil
		default:
		}

		lastCommitTime, err := git.LastCommitUnixtime(localGitTimeout)
		if err != nil {
			notifySend("rp: error finding last commit date")

			return err
		}

		if time.Now().Unix()-lastCommitTime < watchForCommitsInLastSeconds {
			return gitPush(currentBranch)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// git ls-remote --heads origin refs/heads/[branch-name]
func gitPush(currentBranch string) error {
	out, err := grace.RunTimed(
		localGitTimeout,
		nil,
		"git",
		"ls-remote",
		"--heads",
		"origin",
		fmt.Sprintf("refs/heads/%s", currentBranch),
	)
	if err != nil {
		notifySend("rp: failed to check for remote branch")

		return err
	}

	if strings.TrimSpace(out.Combine()) == "" {
		notifySend(fmt.Sprintf("rp: remote branch %s not found, not pushing", currentBranch))

		return nil
	}

	out, err = grace.RunTimed(pushTimeout, nil, "git", "push", "origin", currentBranch)
	if errors.Is(err, grace.ErrFailToKill) {
		notifySend(fmt.Sprintf("rp: timeout reached, but failed to kill (pushing %s)", currentBranch))
	} else if errors.Is(err, grace.ErrTimeout) {
		notifySend(fmt.Sprintf("rp: push timeout: %s", currentBranch))
	} else if err != nil {
		notifySend("rp: error executing git push")
	}

	if err != nil {
		return err
	}

	if out.ExitCode != 0 {
		notifySend(fmt.Sprintf("rp: push failed with %d exit code", out.ExitCode))

		return errors.New(out.Combine())
	}

	notifySend(fmt.Sprintf("rp: pushed %s", currentBranch))

	return nil
}

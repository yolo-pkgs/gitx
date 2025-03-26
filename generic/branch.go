package generic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yolo-pkgs/grace"

	"github.com/yolo-pkgs/gitx/git"
)

const defaultExecTimeout = 10 * time.Second

func FetchCurrent() error {
	if _, err := grace.RunTimed(
		defaultExecTimeout,
		nil,
		"git",
		"fetch",
		"origin",
	); err != nil {
		return fmt.Errorf("failed fetching current branch: %w", err)
	}

	return nil
}

// TODO: handle non-fast-forward error
func FetchDefault(current string) (string, error) {
	// timeout 5 git fetch origin "${BRANCH}":"${BRANCH}"
	defaultBranch, err := git.DefaultBranch(defaultExecTimeout)
	if err != nil {
		return "", fmt.Errorf("failed detecting default branch: %w", err)
	}

	if defaultBranch == current {
		_, err := grace.RunTimed(defaultExecTimeout, nil, "git", "pull")
		if err != nil {
			return "", fmt.Errorf("failed pulling default branch, which is checked out: %w", err)
		}
	} else {
		if _, err := grace.RunTimed(
			defaultExecTimeout,
			nil,
			"git",
			"fetch",
			"origin",
			fmt.Sprintf("%s:%s", defaultBranch, defaultBranch),
		); err != nil {
			return "", fmt.Errorf("failed direct fetch of default branch: %w", err)
		}
	}

	return defaultBranch, nil
}

func FetchCurrentDefault() (string, string, error) {
	current, err := git.CurrentBranch(defaultExecTimeout)
	if err != nil {
		return "", "", fmt.Errorf("failed getting current branch: %w", err)
	}

	if err = FetchCurrent(); err != nil {
		return "", "", fmt.Errorf("failed fetching current branch: %w", err)
	}

	defaultBranch, err := FetchDefault(current)
	if err != nil {
		return "", "", fmt.Errorf("failed fetching default branch: %w", err)
	}

	return current, defaultBranch, nil
}

func LeftRight(leftRef, rightRef string) (int64, int64, error) {
	output, err := grace.RunTimed(
		defaultExecTimeout,
		nil,
		"git",
		"rev-list",
		"--left-right",
		"--count",
		fmt.Sprintf("%s...%s", leftRef, rightRef),
	)
	if err != nil {
		return 0, 0, fmt.Errorf("failed counting left-right: %w", err)
	}

	behindAhead := output.Combine()

	behindAheadF := strings.Fields(behindAhead)
	if len(behindAhead) < 2 {
		return 0, 0, errors.New("behindAhead fields < 2")
	}
	left, err := strconv.ParseInt(behindAheadF[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed parsing left: %w", err)
	}
	right, err := strconv.ParseInt(behindAheadF[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed parsing right: %w", err)
	}
	return left, right, nil
}

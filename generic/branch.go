package generic

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/yolo-pkgs/grace"
)

const defaultExecTimeout = 10 * time.Second

func CurrentBranch() (string, error) {
	output, err := grace.RunTimed(defaultExecTimeout, "git", "branch", "--show-current")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(output), nil
}

func DefaultBranch() (string, error) {
	output, err := grace.RunTimed(defaultExecTimeout, "git", "branch")
	if err != nil {
		return "", fmt.Errorf("failed to run git branch: %w", err)
	}
	fields := strings.Fields(output)

	usualDefaults := []string{"release", "master", "main"}
	candidates := lo.Intersect(fields, usualDefaults)
	if len(candidates) == 0 {
		return "", errors.New("no default branch found")
	}
	if len(candidates) > 1 {
		return "", fmt.Errorf("multiple candidates for default branch found: %v", candidates)
	}

	return candidates[0], nil
}

func FetchDefault(current string) (string, error) {
	// timeout 5 git fetch origin "${BRANCH}":"${BRANCH}"
	defaultBranch, err := DefaultBranch()
	if err != nil {
		return "", fmt.Errorf("failed detecting default branch: %w", err)
	}

	if defaultBranch == current {
		_, err := grace.RunTimed(defaultExecTimeout, "git", "pull")
		if err != nil {
			return "", fmt.Errorf("failed pulling default branch, which is checked out: %w", err)
		}
	} else {
		if _, err := grace.RunTimed(
			defaultExecTimeout,
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

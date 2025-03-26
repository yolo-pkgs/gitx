package git

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/yolo-pkgs/grace"
)

func CurrentBranch(timeout time.Duration) (string, error) {
	output, err := grace.RunTimed(timeout, nil, "git", "branch", "--show-current")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(output.Combine()), nil
}

func DefaultBranch(timeout time.Duration) (string, error) {
	output, err := grace.RunTimed(timeout, nil, "git", "branch")
	if err != nil {
		return "", fmt.Errorf("failed to run git branch: %w", err)
	}
	fields := strings.Fields(output.Combine())

	var chosenDefault string

	orderedDefaults := []string{"develop", "master", "main"}
	for _, candidate := range orderedDefaults {
		if slices.Contains(fields, candidate) {
			chosenDefault = candidate

			break
		}
	}

	if chosenDefault == "" {
		return "", errors.New("no default branch found")
	}

	return chosenDefault, nil
}

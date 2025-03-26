package git

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
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

	usualDefaults := []string{"develop"}
	candidates := lo.Intersect(fields, usualDefaults)
	if len(candidates) == 0 {
		return "", errors.New("no default branch found")
	}
	if len(candidates) > 1 {
		return "", fmt.Errorf("multiple candidates for default branch found: %v", candidates)
	}

	return candidates[0], nil
}

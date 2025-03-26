package status

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/yolo-pkgs/grace"

	"github.com/yolo-pkgs/gitx/generic"
	"github.com/yolo-pkgs/gitx/git"
)

const defaultExecTimeout = 10 * time.Second

func pushTargetExists(branch string) (bool, error) {
	output, err := grace.RunTimed(defaultExecTimeout, nil, "git", "branch", "-r", "--no-color")
	if err != nil {
		return false, fmt.Errorf("failed getting remote branches: %w", err)
	}
	remotes := strings.Fields(output.Combine())

	return slices.Contains(remotes, fmt.Sprintf("origin/%s", branch)), nil
}

func lastCommitSummary() (string, error) {
	unixTime, err := git.LastCommitUnixtime(defaultExecTimeout)
	if err != nil {
		return "", fmt.Errorf("failed getting last commit unix: %w", err)
	}

	commitTime := time.Unix(unixTime, 0)
	ago := time.Since(commitTime)

	return fmt.Sprintf("\nLast commit: %s ago", ago.Round(time.Second)), nil
}

func CoolStatus() (string, error) {
	current, defaultBranch, err := generic.FetchCurrentDefault()
	if err != nil {
		return "", err
	}

	// ahead/behind default branch
	left, right, err := generic.LeftRight("@", defaultBranch)
	if err != nil {
		return "", fmt.Errorf("failed counting left-right: %w", err)
	}
	leftRightDefault := fmt.Sprintf("DEVELOP BRANCH:\nahead %d; behind %d", left, right)

	// ahead/behind @{push}
	var leftRightPushTarget string
	pushExists, err := pushTargetExists(current)
	if err != nil {
		return "", err
	}
	if pushExists {
		left, right, err = generic.LeftRight("@", "@{push}")
		if err != nil {
			return "", fmt.Errorf("failed counting left-right: %w", err)
		}
		leftRightPushTarget = fmt.Sprintf("REMOTE:\nahead %d; behind %d", left, right)
	} else {
		leftRightPushTarget = "REMOTE:\ndoes not exist"
	}

	// check if merged
	mergedRaw, err := grace.RunTimed(defaultExecTimeout, nil, "git", "branch", "--merged", defaultBranch)
	if err != nil {
		return "", fmt.Errorf("failed getting merged status: %w", err)
	}
	merged := strings.Fields(mergedRaw.Combine())

	var mergedMsg string
	if lo.Contains(merged, current) {
		mergedMsg = "Branch already merged."
	} else {
		mergedMsg = "Branch is not merged."
	}

	lastCommitSum, err := lastCommitSummary()
	if err != nil {
		return "", fmt.Errorf("failed getting last commit summary: %w", err)
	}

	return strings.Join([]string{
		mergedMsg,
		leftRightDefault,
		leftRightPushTarget,
		lastCommitSum,
	}, "\n"), nil
}

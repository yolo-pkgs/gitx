package status

import (
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/yolo-pkgs/grace"

	"github.com/yolo-pkgs/gitx/generic"
)

const defaultExecTimeout = 10 * time.Second

func simpleStatus() {
}

func CoolStatus() (string, error) {
	// simple status
	simple, err := grace.RunTimed(defaultExecTimeout, "git", "status", "--show-stash")
	if err != nil {
		return "", fmt.Errorf("failed getting simple status: %w", err)
	}

	current, err := generic.CurrentBranch()
	if err != nil {
		return "", fmt.Errorf("failed getting current branch: %w", err)
	}

	defaultBranch, err := generic.FetchDefault(current)
	if err != nil {
		return "", fmt.Errorf("failed fetching default branch: %w", err)
	}

	// behind/ahead
	behindAhead, err := grace.RunTimed(
		defaultExecTimeout,
		"git",
		"rev-list",
		"--left-right",
		"--count",
		fmt.Sprintf("%s...%s", defaultBranch, defaultBranch),
	)
	if err != nil {
		return "", fmt.Errorf("failed counting left-right: %w", err)
	}
	behindAheadF := strings.Fields(behindAhead)
	behind := behindAheadF[0]
	ahead := behindAheadF[1]
	behindAheadOut := fmt.Sprintf("Behind %s commit; Ahead %s commit", behind, ahead)

	// check if merged
	// git branch --merged | grep "${CURRENT}"
	mergedRaw, err := grace.RunTimed(defaultExecTimeout, "git", "branch", "--merged")
	if err != nil {
		return "", fmt.Errorf("failed getting merged status: %w", err)
	}
	merged := strings.Fields(mergedRaw)

	var mergedMsg string
	if lo.Contains(merged, current) {
		mergedMsg = "Branch already merged."
	} else {
		mergedMsg = "Branch is not merged."
	}

	return strings.Join([]string{simple, behindAheadOut, mergedMsg}, "\n"), nil
}

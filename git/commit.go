package git

import (
	"strconv"
	"strings"
	"time"

	"github.com/yolo-pkgs/grace"
)

func LastCommitUnixtime(timeout time.Duration) (int64, error) {
	output, err := grace.RunTimed(timeout, nil, "git", "log", "-1", "--format=%ct")
	if err != nil {
		return 0, err
	}

	trimmed := strings.TrimSpace(output.Combine())

	timestamp, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return 0, err
	}

	return timestamp, nil
}

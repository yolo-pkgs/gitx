package generic

import (
	"strconv"
	"strings"

	"github.com/yolo-pkgs/grace"
)

func LastCommitUnixtime() (int64, error) {
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

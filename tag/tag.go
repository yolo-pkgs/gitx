package tag

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/yolo-pkgs/grace"
)

const timeout = 5 * time.Second

func create(tag string, prerel bool) error {
	now := time.Now().UTC()
	if prerel {
		tag = tag + `-rc-` + now.Format(`2006.01.02--15.04.05`)
	}

	_, err := grace.RunTimedSh(timeout, fmt.Sprintf(`git tag %s -m "fix: %s"`, tag, tag))
	if err != nil {
		return fmt.Errorf("failed tagging: %w", err)
	}

	slog.Info("tag created", slog.String("tag", tag))

	return nil
}

func Patch(prerelease bool) error {
	_, err := grace.RunTimedSh(timeout, "git fetch --tags")
	if err != nil {
		return fmt.Errorf("fail fetching tags: %w", err)
	}

	output, err := grace.RunTimedSh(timeout, "git tag")
	if err != nil {
		return fmt.Errorf("fail getting tags: %w", err)
	}

	tags := strings.Fields(output)
	if len(tags) == 0 {
		slog.Info("no tags found")
		return create("v0.0.1", prerelease)
	}

	versions := make([]*version.Version, 0)

	for _, tag := range tags {
		v, err := version.NewVersion(tag)
		if err != nil {
			continue
		}
		versions = append(versions, v)
	}

	if len(versions) == 0 {
		slog.Info("no valid go version tags found")
		return create("v0.0.1", prerelease)
	}

	sort.Sort(version.Collection(versions))
	lastVersion := versions[len(versions)-1]

	segments := lastVersion.Segments64()
	if len(segments) != 3 {
		return errors.New("number of segments in last version != 3")
	}

	major := segments[0]
	minor := segments[1]
	patch := segments[2]

	if lastVersion.Prerelease() == "" {
		patch++
	} else if !prerelease {
		patch++
	}

	newTag := fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	return create(newTag, prerelease)
}

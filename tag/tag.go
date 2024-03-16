package tag

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/yolo-pkgs/grace"
)

func runCmd(ctx context.Context, name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	slog.Info("running process", slog.String("cmd", cmd.String()))
	output, err := grace.Spawn(ctx, cmd)
	if err != nil {
		return "", err
	}

	return output.Combine(), nil
}

func publish(tag string) error {
	ctx := context.Background()
	slog.Info("creating tag", slog.String("tag", tag))

	_, err := runCmd(ctx, "git", "tag", tag, "-m", fmt.Sprintf("'fix: %s'", tag))
	if err != nil {
		return fmt.Errorf("failed tagging: %w", err)
	}

	slog.Info("tag created", slog.String("tag", tag))

	return nil
}

func Patch() error {
	ctx := context.Background()

	_, err := runCmd(ctx, "git", "fetch", "--tags")
	if err != nil {
		return fmt.Errorf("fail fetching tags: %w", err)
	}

	output, err := runCmd(ctx, "git", "tag")
	if err != nil {
		return fmt.Errorf("fail getting tags: %w", err)
	}
	tags := strings.Fields(output)
	if len(tags) == 0 {
		slog.Info("no tags found")
		return publish("v0.0.1")
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
		return publish("v0.0.1")
	}
	sort.Sort(version.Collection(versions))
	lastVersion := versions[len(versions)-1]
	segments := lastVersion.Segments64()
	if len(segments) != 3 {
		return errors.New("number of segments in last version != 3")
	}

	newTag := fmt.Sprintf("v%d.%d.%d", segments[0], segments[1], segments[2]+1)
	return publish(newTag)
}

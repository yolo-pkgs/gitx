package tag

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/yolo-pkgs/grace"
)

const timeout = 5 * time.Second

type Version struct {
	Major int64
	Minor int64
	Patch int64
}

func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func lastReleaseTag(versions []Version) Version {
	majors := lo.Map(versions, func(item Version, _ int) int64 { return item.Major })
	maxMajor := lo.Max(majors)
	tagsWithMaxMajor := lo.Filter(versions, func(item Version, _ int) bool { return item.Major == maxMajor })

	minors := lo.Map(tagsWithMaxMajor, func(item Version, _ int) int64 { return item.Minor })
	maxMinor := lo.Max(minors)
	tagsWithMaxMinor := lo.Filter(tagsWithMaxMajor, func(item Version, _ int) bool { return item.Minor == maxMinor })

	patches := lo.Map(tagsWithMaxMinor, func(item Version, _ int) int64 { return item.Patch })
	maxPatch := lo.Max(patches)

	return Version{
		Major: maxMajor,
		Minor: maxMinor,
		Patch: maxPatch,
	}
}

func parseTag(tag string) (Version, bool) {
	tag = strings.TrimPrefix(tag, "v")
	fields := strings.Split(tag, ".")
	if len(fields) != 3 {
		return Version{}, false
	}

	major, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return Version{}, false
	}

	minor, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return Version{}, false
	}

	patch, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return Version{}, false
	}

	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, true
}

func create(tag string, prerel bool) error {
	now := time.Now().UTC()
	if prerel {
		tag = tag + `-rc-` + now.Format(`2006.01.02--15.04.05`)
	}

	// _, err := grace.RunTimedSh(timeout, fmt.Sprintf(`git tag %s -m "fix: %s"`, tag, tag))
	// if err != nil {
	// 	return fmt.Errorf("failed tagging: %w", err)
	// }

	fmt.Printf("tag created: %s\n", tag)

	return nil
}

func Patch(prerelease bool) error {
	r, err := regexp.Compile(`^v\d+\.\d+\.\d+$`)
	if err != nil {
		return fmt.Errorf("failed to compile release tag regex: %w", err)
	}

	_, err = grace.RunTimedSh(timeout, "git fetch --tags")
	if err != nil {
		return fmt.Errorf("fail fetching tags: %w", err)
	}

	output, err := grace.RunTimedSh(timeout, "git tag")
	if err != nil {
		return fmt.Errorf("fail getting tags: %w", err)
	}

	tags := strings.Fields(output)
	if len(tags) == 0 {
		fmt.Println("no tags found")
		return create("v0.0.1", prerelease)
	}

	releaseTags := make([]Version, 0)
	for _, tag := range tags {
		if !r.MatchString(tag) {
			continue
		}

		parsed, ok := parseTag(tag)
		if !ok {
			continue
		}

		releaseTags = append(releaseTags, parsed)
	}

	if len(releaseTags) == 0 {
		fmt.Println("no valid go version tags found")
		return create("v0.0.1", prerelease)
	}

	lastRel := lastReleaseTag(releaseTags)
	fmt.Printf("last release tag: %s\n", lastRel.String())

	newTag := fmt.Sprintf("v%d.%d.%d", lastRel.Major, lastRel.Minor, lastRel.Patch+1)
	return create(newTag, prerelease)
}

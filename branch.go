package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/yolo-pkgs/grace"

	"github.com/yolo-pkgs/gitx/generic"
	"github.com/yolo-pkgs/gitx/system"
)

const (
	branchPrefix       = "s"
	gitxBranchFilePath = "sys/.gitx_branch"
)

func fromDefaultBranch() (bool, error) {
	current, err := generic.CurrentBranch()
	if err != nil {
		return false, err
	}

	def, err := generic.DefaultBranch()
	if err != nil {
		return false, err
	}

	if current != def {
		return false, nil
	}

	return true, nil
}

func fetchAll() error {
	_, err := grace.RunTimed(defaultExecTimeout, "git", "fetch", "--all")
	return err
}

func listBranches() error {
	if err := fetchAll(); err != nil {
		return err
	}

	output, err := grace.RunTimed(defaultExecTimeout, "git", "branch", "--all")
	if err != nil {
		return err
	}

	fmt.Println(output)
	return nil
}

func globalBranchID() (int64, error) {
	home, err := system.GetHome()
	if err != nil {
		return 0, err
	}

	raw, err := os.ReadFile(fmt.Sprintf("%s/%s", home, gitxBranchFilePath))
	if err != nil {
		return 0, fmt.Errorf("failed to read .gitx_branch: %w", err)
	}

	data := strings.TrimSpace(string(raw))
	gid, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse global branch id: %w", err)
	}
	return gid, nil
}

func writeNewBranchGID() (int64, error) {
	home, err := system.GetHome()
	if err != nil {
		return 0, err
	}

	gid, err := globalBranchID()
	if err != nil {
		return 0, err
	}
	newGID := gid + 1

	f, err := os.OpenFile(fmt.Sprintf("%s/%s", home, gitxBranchFilePath), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return 0, fmt.Errorf("failed to open %s: %w", gitxBranchFilePath, err)
	}

	_, err = io.WriteString(f, strconv.FormatInt(newGID, 10))
	if err != nil {
		return 0, fmt.Errorf("failed to write new gid: %w", err)
	}

	if err := f.Close(); err != nil {
		return 0, fmt.Errorf("failed to close %s: %w", gitxBranchFilePath, err)
	}

	return newGID, nil
}

func randomWord() (string, error) {
	home, err := system.GetHome()
	if err != nil {
		return "", err
	}

	// TODO: wordlist to reliable place
	output, err := grace.RunTimed(
		defaultExecTimeout,
		"shuf",
		"-n",
		"1",
		fmt.Sprintf("%s/dev/dotfiles/config/nvim/100k.txt", home),
	)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

func getModifier(fromDefault, _ bool) string {
	if fromDefault {
		return ""
	}
	return "x"
}

func makeSourceMark(modifier string) (string, error) {
	if modifier == "f" || modifier == "x" {
		current, err := generic.CurrentBranch()
		if err != nil {
			return "", err
		}
		currentFields := strings.Split(current, "-")
		last := currentFields[len(currentFields)-1]
		return "-" + last, nil
	}
	return "", nil
}

func createRandomBranch(gid int64, fromDefault, xMark bool) error {
	modifier := getModifier(fromDefault, xMark)
	// sourceMark, err := makeSourceMark(modifier)
	// if err != nil {
	// 	return err
	// }

	// randomWord, err := randomWord()
	// if err != nil {
	// 	return fmt.Errorf("failed to generate random word: %w", err)
	// }

	branchName := fmt.Sprintf("%s%d%s", branchPrefix, gid, modifier)
	_, err := grace.RunTimed(defaultExecTimeout, "git", "checkout", "-b", branchName)
	return err
}

func createGlobalBranch(gid int64, name string, fromDefault, xMark bool) error {
	modifier := getModifier(fromDefault, xMark)
	// sourceMark, err := makeSourceMark(modifier)
	// if err != nil {
	// 	return err
	// }

	branchName := fmt.Sprintf("%s%d%s-%s", branchPrefix, gid, modifier, name)
	_, err := grace.RunTimed(defaultExecTimeout, "git", "checkout", "-b", branchName)
	return err
}

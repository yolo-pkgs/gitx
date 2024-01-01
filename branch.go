package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/yolo-pkgs/grace"
)

const (
	branchPrefix       = "samira/"
	gitxBranchFilePath = "sys/.gitx_branch"
)

func getHome() (string, error) {
	home := os.Getenv("")
	if home == "" {
		return "", errors.New("HOME env not found")
	}
	return home, nil
}

func currentBranch() (string, error) {
	output, err := grace.RunTimed(defaultExecTimeout, "git", "branch", "--show-current")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(output), nil
}

func defaultBranch() (string, error) {
	output, err := grace.RunTimed(defaultExecTimeout, "git", "branch")
	if err != nil {
		return "", fmt.Errorf("failed to run git branch: %w", err)
	}
	fields := strings.Fields(output)

	usualDefaults := []string{"release", "master", "main"}
	candidates := lo.Intersect(fields, usualDefaults)
	if len(candidates) == 0 {
		return "", errors.New("no default branch found")
	}
	if len(candidates) > 1 {
		return "", fmt.Errorf("multiple candidates for default branch found: %v", candidates)
	}

	return candidates[0], nil
}

func onlyFromDefaultBranch() error {
	current, err := currentBranch()
	if err != nil {
		return err
	}

	def, err := defaultBranch()
	if err != nil {
		return err
	}

	if current != def {
		return errors.New("branch actions are only allowed from default branch")
	}

	return nil
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
	home, err := getHome()
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
	home, err := getHome()
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
	home, err := getHome()
	if err != nil {
		return "", err
	}

	// TODO: wordlist to reliable place
	output, err := grace.RunTimed(defaultExecTimeout, "shuf", "-n", "1", fmt.Sprintf("%s/dev/dotfiles/config/nvim/100k.txt", home))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

func createRandomBranch(gid int64) error {
	randomWord, err := randomWord()
	if err != nil {
		return fmt.Errorf("failed to generate random word: %w", err)
	}

	branchName := fmt.Sprintf("%sg%d-%s", branchPrefix, gid, randomWord)
	_, err = grace.RunTimed(defaultExecTimeout, "git", "checkout", "-b", branchName)
	return err
}

func createGlobalBranch(gid int64, name string) error {
	branchName := fmt.Sprintf("%sg%d-%s", branchPrefix, gid, name)
	_, err := grace.RunTimed(defaultExecTimeout, "git", "checkout", "-b", branchName)
	return err
}

package system

import (
	"fmt"
	"os"
)

func GetHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("HOME env error: %w", err)
	}
	return home, nil
}

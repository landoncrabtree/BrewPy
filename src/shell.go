package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func updateShellProfile() error {
	config := loadConfig()
	
	content, err := os.ReadFile(config.ShellRC)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lines := strings.Split(string(content), "\n")

	// Remove existing brewpy init block
	startIdx, endIdx := -1, -1
	for i, line := range lines {
		if line == initComment {
			startIdx = i
		}
		if line == initEndComment {
			endIdx = i
		}
	}

	if startIdx != -1 && endIdx != -1 && startIdx < endIdx {
		lines = append(lines[:startIdx], lines[endIdx+1:]...)
	}

	// Add brewpy init block
	initBlock := []string{
		"",
		initComment,
		`eval "$(brewpy init)"`,
		initEndComment,
	}

	// Append init block at the end
	lines = append(lines, initBlock...)

	// Write back to file
	return os.WriteFile(config.ShellRC, []byte(strings.Join(lines, "\n")), fs.ModePerm)
}

func outputShellInit(homeDir string) {
	config := loadConfig()
	shimsPath := getShimsDir(config.BrewPyDir)
	fmt.Printf("export PATH=\"%s:$PATH\"", shimsPath)
} 
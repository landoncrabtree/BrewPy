package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func createSymlinks(version string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	shimsPath := filepath.Join(homeDir, brewpyDir, shimsDir)

	// Create ~/.brewpy/shims directory if it doesn't exist
	err = os.MkdirAll(shimsPath, 0755)
	if err != nil {
		return err
	}

	// Extract version number from "Python3.11" -> "3.11"
	ver := strings.TrimPrefix(version, "Python")
	binDir := getBinDir()

	// Define the symlinks to create
	links := map[string]string{
		"python":  filepath.Join(binDir, "python"+ver),
		"python3": filepath.Join(binDir, "python"+ver),
		"pip":     filepath.Join(binDir, "pip"+ver),
		"pip3":    filepath.Join(binDir, "pip"+ver),
	}

	// Remove existing symlinks and create new ones
	for linkName, target := range links {
		linkPath := filepath.Join(shimsPath, linkName)

		// Remove existing symlink if it exists
		if _, err := os.Lstat(linkPath); err == nil {
			os.Remove(linkPath)
		}

		// Create new symlink
		err = os.Symlink(target, linkPath)
		if err != nil {
			return fmt.Errorf("failed to create symlink %s -> %s: %w", linkPath, target, err)
		}
	}

	return nil
} 
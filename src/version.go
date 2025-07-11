package main

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
)

func getBinDir() string {
	if runtime.GOARCH == "arm64" {
		return "/opt/homebrew/bin"
	}
	return "/usr/local/bin"
}

func findPythonVersions() ([]string, error) {
	binDir := getBinDir()
	files, err := os.ReadDir(binDir)
	if err != nil {
		return nil, err
	}

	// regex to match python3.11 or python3.12 etc
	pythonRe := regexp.MustCompile(`^python(\d+\.\d+)$`)

	versionMap := map[string]struct{}{}

	for _, file := range files {
		name := file.Name()
		matches := pythonRe.FindStringSubmatch(name)
		if len(matches) == 2 {
			ver := "Python" + matches[1]
			versionMap[ver] = struct{}{}
		}
	}

	versions := make([]string, 0, len(versionMap))
	for v := range versionMap {
		versions = append(versions, v)
	}
	sort.Strings(versions)

	return versions, nil
}

func getCurrentVersion() string {
	config := loadConfig()
	
	shimsDir := getShimsDir(config.BrewPyDir)
	pythonShim := filepath.Join(shimsDir, "python")
	if _, err := os.Lstat(pythonShim); os.IsNotExist(err) {
		return ""
	}

	target, err := os.Readlink(pythonShim)
	if err != nil {
		return ""
	}

	// Extract version from target path like /opt/homebrew/bin/python3.11
	base := filepath.Base(target)
	re := regexp.MustCompile(`python(\d+\.\d+)`)
	matches := re.FindStringSubmatch(base)
	if len(matches) == 2 {
		return "Python" + matches[1]
	}

	return ""
} 
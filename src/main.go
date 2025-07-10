package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "versions":
		handleVersions()
	case "use":
		handleUse()
	case "init":
		handleInit()
	case "current":
		handleCurrent()
	case "--help", "-h", "help":
		showUsage()
	default:
		showUsage()
	}
}

func handleVersions() {
	displayVersionsHeader()
	
	versions, err := findPythonVersions()
	if err != nil {
		log.Fatal(red("Error finding Python versions: "), err)
	}
	
	if len(versions) == 0 {
		fmt.Printf("%s\n", yellow("No Python versions found. Install Python via Homebrew first."))
		return
	}
	
	current := getCurrentVersion()
	displayVersionsList(versions, current)
}

func handleUse() {
	versions, err := findPythonVersions()
	if err != nil {
		log.Fatal(red("Error finding Python versions: "), err)
	}
	
	if len(versions) == 0 {
		fmt.Printf("%s\n", red("No Python versions found. Install Python via Homebrew first."))
		return
	}
	
	var version string
	if len(os.Args) >= 3 {
		version = os.Args[2]
	} else {
		version, err = promptSelectVersion(versions)
		if err != nil {
			log.Fatal(red("Error selecting version: "), err)
		}
	}
	
	if !contains(versions, version) {
		log.Fatalf("%s %s", red("Version not found:"), version)
	}
	
	err = createSymlinks(version)
	if err != nil {
		log.Fatal(red("Error creating symlinks: "), err)
	}
	
	err = updateShellProfile()
	if err != nil {
		log.Fatal(red("Error updating shell profile: "), err)
	}
	
	displaySuccessMessage(version)
}

func handleInit() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(red("Error getting home directory: "), err)
	}
	
	outputShellInit(homeDir)
}

func handleCurrent() {
	current := getCurrentVersion()
	displayCurrentVersion(current)
}

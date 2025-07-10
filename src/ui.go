package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

var (
	// Color functions for better UI
	green   = color.New(color.FgGreen).SprintFunc()
	blue    = color.New(color.FgBlue).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	cyan    = color.New(color.FgCyan).SprintFunc()
	bold    = color.New(color.Bold).SprintFunc()
)

func showUsage() {
	fmt.Printf(`%s
  %s - list installed python versions
  %s - set python version (e.g. Python3.11). If no version given, prompts selection
  %s - output shell configuration
  %s - show currently active python version
`,
		bold("🍺 BrewPy - Python Version Manager"),
		cyan("brewpy versions"),
		cyan("brewpy use [version]"),
		cyan("brewpy init"),
		cyan("brewpy current"),
	)
}

func displayVersionsHeader() {
	fmt.Printf("%s\n", bold("🔍 Available Python Versions:"))
}

func displayVersionsList(versions []string, current string) {
	for _, v := range versions {
		if v == current {
			fmt.Printf("  %s %s\n", green("●"), green(v))
		} else {
			fmt.Printf("  %s %s\n", "○", v)
		}
	}
}

func displayCurrentVersion(current string) {
	if current == "" {
		fmt.Printf("%s\n", yellow("No Python version currently managed by BrewPy"))
	} else {
		fmt.Printf("%s %s\n", green("Current Python version:"), green(current))
	}
}

func displaySuccessMessage(version string) {
	fmt.Printf("%s %s\n", green("✓ Successfully switched to"), green(version))
	fmt.Printf("%s\n", yellow("Restart your terminal or run 'source ~/.zshrc' to apply changes."))
}

func promptSelectVersion(versions []string) (string, error) {
	prompt := promptui.Select{
		Label: fmt.Sprintf("%s Select Python Version", "🐍"),
		Items: versions,
		Size:  10,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   fmt.Sprintf("%s {{ . | cyan }}", "▸"),
			Inactive: "  {{ . }}",
			Selected: fmt.Sprintf("%s {{ . | green }}", "✓"),
		},
	}

	_, result, err := prompt.Run()
	return result, err
} 
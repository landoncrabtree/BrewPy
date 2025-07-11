package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

type Config struct {
	ShellRC   string `json:"shell_rc"`
	BrewPyDir string `json:"brewpy_dir"`
}

// getDefaultBrewPyDir returns the default BrewPy directory
func getDefaultBrewPyDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".brewpy")
}

// getShimsDir returns the shims directory path based on BrewPyDir
func getShimsDir(brewPyDir string) string {
	return filepath.Join(brewPyDir, "shims")
}

// getConfigPath returns the path to the config file based on BrewPyDir
func getConfigPath(brewPyDir string) string {
	return filepath.Join(brewPyDir, "config.json")
}

// findConfigFile attempts to find the config file, checking current BrewPyDir first, then default location
func findConfigFile(brewPyDir string) (string, bool) {
	// First, try the specified BrewPyDir
	configPath := getConfigPath(brewPyDir)
	if _, err := os.Stat(configPath); err == nil {
		return configPath, true
	}

	// If not found and we're not looking in the default location, try default
	defaultBrewPyDir := getDefaultBrewPyDir()
	if brewPyDir != defaultBrewPyDir {
		defaultConfigPath := getConfigPath(defaultBrewPyDir)
		if _, err := os.Stat(defaultConfigPath); err == nil {
			return defaultConfigPath, true
		}
	}

	return configPath, false
}

func getDefaultConfig() Config {
	brewPyDir := getDefaultBrewPyDir()
	
	return Config{
		BrewPyDir: brewPyDir,
		ShellRC:   detectShellRC(),
	}
}

func detectShellRC() string {
	homeDir, _ := os.UserHomeDir()
	
	// Check for common shell RC files in order of preference
	rcFiles := []string{".zshrc", ".bashrc", ".config/fish/config.fish"}
	
	for _, rcFile := range rcFiles {
		fullPath := filepath.Join(homeDir, rcFile)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}
	
	// Default to .zshrc if none found
	return filepath.Join(homeDir, ".zshrc")
}

// initConfig creates the BrewPy directory and initializes config if needed
func initConfig(config Config) error {
	// Create BrewPy directory
	if err := os.MkdirAll(config.BrewPyDir, 0755); err != nil {
		return fmt.Errorf("failed to create BrewPy directory: %w", err)
	}
	
	// Create shims directory
	shimsDir := getShimsDir(config.BrewPyDir)
	if err := os.MkdirAll(shimsDir, 0755); err != nil {
		return fmt.Errorf("failed to create shims directory: %w", err)
	}
	
	// Save config to ensure it exists in the correct location
	return saveConfig(config)
}

func loadConfig() Config {
	// Start with default config
	config := getDefaultConfig()
	
	// Try to find existing config file
	configPath, exists := findConfigFile(config.BrewPyDir)
	
	if !exists {
		// No config file found, initialize with defaults
		if err := initConfig(config); err != nil {
			fmt.Printf("%s Failed to initialize config: %v\n", yellow("Warning:"), err)
		}
		return config
	}
	
	// Read existing config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("%s Failed to read config file, using defaults: %v\n", yellow("Warning:"), err)
		return config
	}
	
	// Parse JSON
	var loadedConfig Config
	if err := json.Unmarshal(data, &loadedConfig); err != nil {
		fmt.Printf("%s Failed to parse config file, using defaults: %v\n", yellow("Warning:"), err)
		return config
	}
	
	// Migrate config if BrewPyDir has changed and we loaded from default location
	if configPath == getConfigPath(getDefaultBrewPyDir()) && loadedConfig.BrewPyDir != getDefaultBrewPyDir() {
		if err := migrateConfig(loadedConfig, configPath); err != nil {
			fmt.Printf("%s Failed to migrate config: %v\n", yellow("Warning:"), err)
		}
	}
	
	return loadedConfig
}

// migrateConfig moves the config file to the correct location based on BrewPyDir
func migrateConfig(config Config, oldConfigPath string) error {
	// Initialize the new location
	if err := initConfig(config); err != nil {
		return err
	}
	
	// Remove old config file if it's in the default location and we've moved
	if oldConfigPath == getConfigPath(getDefaultBrewPyDir()) && config.BrewPyDir != getDefaultBrewPyDir() {
		os.Remove(oldConfigPath)
	}
	
	return nil
}

func saveConfig(config Config) error {
	configPath := getConfigPath(config.BrewPyDir)
	
	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Marshal to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	// Write config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

func handleConfigure() {
	config := loadConfig()
	
	fmt.Printf("%s\n", bold("🔧 BrewPy Configuration"))
	fmt.Printf("Configure BrewPy settings interactively.\n\n")
	
	// Show current configuration
	fmt.Printf("%s Current Configuration:\n", bold("📋"))
	fmt.Printf("  BrewPy directory: %s\n", blue(config.BrewPyDir))
	fmt.Printf("  Shims directory:  %s\n", blue(getShimsDir(config.BrewPyDir)))
	fmt.Printf("  Shell RC file:    %s\n", blue(config.ShellRC))
	fmt.Printf("\n")
	
	// Ask what to configure
	configChoice, err := promptConfigChoice()
	if err != nil {
		fmt.Printf("%s %v\n", red("Error:"), err)
		return
	}
	
	switch configChoice {
	case "brewpy_dir":
		if err := configureBrewPyDirectory(&config); err != nil {
			fmt.Printf("%s %v\n", red("Error:"), err)
			return
		}
		
	case "shell_rc":
		if err := configureShellRC(&config); err != nil {
			fmt.Printf("%s %v\n", red("Error:"), err)
			return
		}
		
	case "all":
		if err := configureAll(&config); err != nil {
			fmt.Printf("%s %v\n", red("Error:"), err)
			return
		}
		
	case "reset":
		if confirmed, _ := promptConfirmReset(); confirmed {
			config = getDefaultConfig()
		} else {
			fmt.Printf("%s Configuration reset cancelled.\n", yellow("Cancelled"))
			return
		}
		
	default:
		fmt.Printf("%s No changes made.\n", yellow("Cancelled"))
		return
	}
	
	// Save and initialize configuration
	if err := saveConfig(config); err != nil {
		fmt.Printf("%s Failed to save configuration: %v\n", red("Error:"), err)
		return
	}
	
	if err := initConfig(config); err != nil {
		fmt.Printf("%s Failed to initialize configuration: %v\n", red("Error:"), err)
		return
	}
	
	fmt.Printf("\n%s Configuration saved successfully!\n", green("✓"))
	fmt.Printf("Configuration file: %s\n", cyan(getConfigPath(config.BrewPyDir)))
	fmt.Printf("\n%s Updated settings:\n", bold("📋"))
	fmt.Printf("  BrewPy directory: %s\n", config.BrewPyDir)
	fmt.Printf("  Shims directory:  %s\n", getShimsDir(config.BrewPyDir))
	fmt.Printf("  Shell RC file:    %s\n", config.ShellRC)
	
	fmt.Printf("\n%s Run 'brewpy use' to apply changes to your Python setup.\n", yellow("Note:"))
}

func configureBrewPyDirectory(config *Config) error {
	newDir, err := promptBrewPyDirectory(config.BrewPyDir)
	if err != nil {
		return err
	}
	if newDir != "" {
		config.BrewPyDir = expandPath(newDir)
	}
	return nil
}

func configureShellRC(config *Config) error {
	newRC, err := promptShellRC(config.ShellRC)
	if err != nil {
		return err
	}
	if newRC != "" {
		config.ShellRC = newRC
	}
	return nil
}

func configureAll(config *Config) error {
	if err := configureBrewPyDirectory(config); err != nil {
		return err
	}
	return configureShellRC(config)
}

func promptConfigChoice() (string, error) {
	items := []string{
		"BrewPy directory (where config and shims are stored)",
		"Shell RC file (where 'brewpy init' will be added)",
		"Configure all settings",
		"Reset to defaults",
		"Cancel",
	}
	
	prompt := promptui.Select{
		Label: "🔧 What would you like to configure?",
		Items: items,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   fmt.Sprintf("%s {{ . | cyan }}", "▸"),
			Inactive: "  {{ . }}",
			Selected: fmt.Sprintf("%s {{ . | green }}", "✓"),
		},
		Size: 5,
	}
	
	index, _, err := prompt.Run()
	if err != nil {
		return "", err
	}
	
	choices := []string{"brewpy_dir", "shell_rc", "all", "reset", "cancel"}
	return choices[index], nil
}

func promptBrewPyDirectory(current string) (string, error) {
	validate := func(input string) error {
		if input == "" {
			return nil // Empty is allowed (keep current)
		}
		
		expandedPath := expandPath(input)
		parentDir := filepath.Dir(expandedPath)
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			return fmt.Errorf("parent directory %s does not exist", parentDir)
		}
		return nil
	}
	
	prompt := promptui.Prompt{
		Label:   "🔧 BrewPy directory (stores config and shims)",
		Default: current,
		Validate: validate,
		Templates: &promptui.PromptTemplates{
			Prompt:  "{{ . | cyan }}: ",
			Valid:   "{{ . | cyan }}: ",
			Invalid: "{{ . | red }}: ",
			Success: "{{ . | green }}: ",
		},
	}
	
	result, err := prompt.Run()
	return strings.TrimSpace(result), err
}

func promptShellRC(currentValue string) (string, error) {
	homeDir, _ := os.UserHomeDir()
	
	// Build list of shell RC options
	type shellRCOption struct {
		Name   string
		Path   string
		Exists bool
	}
	
	options := []shellRCOption{
		{"Keep current", currentValue, true},
		{"Enter custom path", "", false},
	}
	
	// Add common shell RC files
	rcFiles := []string{".zshrc", ".bashrc", ".bash_profile", ".profile", ".dashrc", ".config/fish/config.fish"}
	for _, rcFile := range rcFiles {
		fullPath := filepath.Join(homeDir, rcFile)
		exists := false
		if _, err := os.Stat(fullPath); err == nil {
			exists = true
		}
		
		// Don't add if it's already the current value
		if fullPath != currentValue {
			options = append(options, shellRCOption{
				Name:   rcFile,
				Path:   fullPath,
				Exists: exists,
			})
		}
	}
	
	// Create display items
	items := make([]string, len(options))
	for i, option := range options {
		if option.Name == "Keep current" {
			items[i] = fmt.Sprintf("Keep current (%s)", filepath.Base(option.Path))
		} else if option.Name == "Enter custom path" {
			items[i] = "Enter custom path..."
		} else {
			status := "not found"
			if option.Exists {
				status = "exists"
			}
			items[i] = fmt.Sprintf("%s (%s)", option.Name, status)
		}
	}
	
	prompt := promptui.Select{
		Label: "🐚 Select shell RC file",
		Items: items,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   fmt.Sprintf("%s {{ . | cyan }}", "▸"),
			Inactive: "  {{ . }}",
			Selected: fmt.Sprintf("%s {{ . | green }}", "✓"),
		},
		Size: 10,
	}
	
	index, _, err := prompt.Run()
	if err != nil {
		return "", err
	}
	
	selected := options[index]
	
	// Handle custom path option
	if selected.Name == "Enter custom path" {
		validate := func(input string) error {
			if input == "" {
				return fmt.Errorf("path cannot be empty")
			}
			expandedPath := expandPath(input)
			parentDir := filepath.Dir(expandedPath)
			if _, err := os.Stat(parentDir); os.IsNotExist(err) {
				return fmt.Errorf("parent directory %s does not exist", parentDir)
			}
			return nil
		}
		
		customPrompt := promptui.Prompt{
			Label:    "🐚 Shell RC file path",
			Default:  currentValue,
			Validate: validate,
			Templates: &promptui.PromptTemplates{
				Prompt:  "{{ . | cyan }}: ",
				Valid:   "{{ . | cyan }}: ",
				Invalid: "{{ . | red }}: ",
				Success: "{{ . | green }}: ",
			},
		}
		
		result, err := customPrompt.Run()
		return expandPath(strings.TrimSpace(result)), err
	}
	
	// Handle keep current option
	if selected.Name == "Keep current" {
		return "", nil // Empty string means keep current
	}
	
	return selected.Path, nil
}

func promptConfirmReset() (bool, error) {
	prompt := promptui.Prompt{
		Label:     "⚠️  Reset all settings to defaults? (y/N)",
		Default:   "N",
		AllowEdit: true,
	}
	
	result, err := prompt.Run()
	if err != nil {
		return false, err
	}
	
	return strings.ToLower(strings.TrimSpace(result)) == "y", nil
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

func handleConfigShow() {
	config := loadConfig()
	
	fmt.Printf("%s\n", bold("📋 BrewPy Configuration"))
	fmt.Printf("Configuration file: %s\n\n", cyan(getConfigPath(config.BrewPyDir)))
	
	fmt.Printf("BrewPy directory: %s\n", config.BrewPyDir)
	fmt.Printf("Shims directory:  %s\n", getShimsDir(config.BrewPyDir))
	fmt.Printf("Shell RC file:    %s\n", config.ShellRC)
	
	fmt.Printf("\n%s Status:\n", bold("📊"))
	
	// Check if directories exist
	if _, err := os.Stat(config.BrewPyDir); os.IsNotExist(err) {
		fmt.Printf("  %s BrewPy directory does not exist\n", yellow("⚠"))
	} else {
		fmt.Printf("  %s BrewPy directory exists\n", green("✓"))
	}
	
	shimsDir := getShimsDir(config.BrewPyDir)
	if _, err := os.Stat(shimsDir); os.IsNotExist(err) {
		fmt.Printf("  %s Shims directory does not exist\n", yellow("⚠"))
	} else {
		fmt.Printf("  %s Shims directory exists\n", green("✓"))
	}
	
	if _, err := os.Stat(config.ShellRC); os.IsNotExist(err) {
		fmt.Printf("  %s Shell RC file does not exist\n", yellow("⚠"))
	} else {
		fmt.Printf("  %s Shell RC file exists\n", green("✓"))
	}
} 
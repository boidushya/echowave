package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Dependency represents an external tool required for EchoWave operation.
// Contains command name, executable path, and platform-specific installation instructions.
type Dependency struct {
	Name        string
	Command     string
	InstallDocs map[string]string
}

var dependencies = []Dependency{
	{
		Name:    "ffmpeg",
		Command: "ffmpeg",
		InstallDocs: map[string]string{
			"darwin":  "brew install ffmpeg",
			"linux":   "sudo apt-get install ffmpeg  # Ubuntu/Debian\nsudo yum install ffmpeg     # CentOS/RHEL",
			"windows": "Download from https://ffmpeg.org/download.html",
		},
	},
	{
		Name:    "openai-whisper",
		Command: "whisper",
		InstallDocs: map[string]string{
			"darwin":  "pip install openai-whisper",
			"linux":   "pip install openai-whisper",
			"windows": "pip install openai-whisper",
		},
	},
	{
		Name:    "yt-dlp",
		Command: "yt-dlp",
		InstallDocs: map[string]string{
			"darwin":  "brew install yt-dlp\n# OR\npip install yt-dlp",
			"linux":   "pip install yt-dlp\n# OR\nsudo apt-get install yt-dlp  # Ubuntu 22.04+",
			"windows": "pip install yt-dlp\n# OR download from https://github.com/yt-dlp/yt-dlp/releases",
		},
	},
}

// checkDependency verifies if a specific dependency is installed and available in the system PATH.
// It takes a Dependency struct and attempts to locate the corresponding command using exec.LookPath.
// Returns true if the dependency is found and executable, false otherwise. This function is used
// to validate individual runtime requirements before proceeding with audio processing operations.
func checkDependency(dep Dependency) bool {
	_, err := exec.LookPath(dep.Command)
	return err == nil
}

// showInstallInstructions displays platform-specific installation commands for a missing dependency.
// It takes a Dependency struct and prints formatted installation instructions based on the current
// operating system (runtime.GOOS). The output includes colored formatting with comments in secondary
// color and commands in white. If no platform-specific instructions exist, it shows a generic
// installation message. This function helps users resolve missing dependencies quickly.
func showInstallInstructions(dep Dependency) {
	subheader("Installing " + dep.Name)

	if instructions, exists := dep.InstallDocs[runtime.GOOS]; exists {
		lines := strings.Split(instructions, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "#") {
				fmt.Printf("%s%s\n", prefix(), colorize(line, SecondaryColor))
			} else {
				fmt.Printf("%s%s\n", prefix(), colorize(line, White))
			}
		}
	} else {
		fmt.Printf("%s%s\n", prefix(), colorize("Please install "+dep.Name+" for your operating system", MutedColor))
	}
	fmt.Println()
}

// checkAllDependencies validates that all required external tools are installed and accessible.
// It iterates through the global dependencies slice, checking each one using checkDependency.
// For each dependency, it prints a success or error message with appropriate formatting.
// If any dependencies are missing, it displays comprehensive installation instructions for
// the current platform and returns false. Returns true only if all dependencies are satisfied,
// allowing the main program to proceed with audio processing operations.
func checkAllDependencies() bool {
	step("Checking dependencies...")

	allPresent := true
	var missing []Dependency

	for _, dep := range dependencies {
		if checkDependency(dep) {
			success(dep.Name + " found")
		} else {
			errorMsg(dep.Name + " not found")
			missing = append(missing, dep)
			allPresent = false
		}
	}

	if !allPresent {
		fmt.Println()
		warning("Missing dependencies: " + strings.Join(getMissingNames(missing), ", "))
		fmt.Println()
		header("Installation Instructions")
		for _, dep := range missing {
			showInstallInstructions(dep)
		}
		fmt.Println()
		info("After installing the missing dependencies, please run the command again.")
		return false
	}

	success("All dependencies are installed!")
	return true
}

// getMissingNames extracts the names of missing dependencies from a slice of Dependency structs.
// It takes a slice of Dependency objects and returns a slice of strings containing just the
// name field from each dependency. This utility function is used to create readable lists
// of missing dependencies for error messages and user feedback during dependency checking.
func getMissingNames(deps []Dependency) []string {
	names := make([]string, len(deps))
	for i, dep := range deps {
		names[i] = dep.Name
	}
	return names
}

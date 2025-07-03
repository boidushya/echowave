package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

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

func checkDependency(dep Dependency) bool {
	_, err := exec.LookPath(dep.Command)
	return err == nil
}

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

func getMissingNames(deps []Dependency) []string {
	names := make([]string, len(deps))
	for i, dep := range deps {
		names[i] = dep.Name
	}
	return names
}
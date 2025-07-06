package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var VERSION = "dev"

const (
	GITHUB_REPO          = "boidushya/echowave"
	GITHUB_API_URL       = "https://api.github.com/repos/" + GITHUB_REPO + "/releases/latest"
	UPDATE_CHECK_TIMEOUT = 3 * time.Second
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func getCurrentVersion() string {
	if VERSION != "dev" {
		return VERSION
	}

	if version := getVersionFromGit(); version != "" {
		return version
	}
	return VERSION
}

func getVersionFromGit() string {
	cmd := exec.Command("git", "describe", "--tags", "--always", "--dirty")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	version := strings.TrimSpace(string(output))
	return strings.TrimPrefix(version, "v")
}

func getLatestVersion() (string, error) {
	client := &http.Client{Timeout: UPDATE_CHECK_TIMEOUT}
	resp, err := client.Get(GITHUB_API_URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch latest version: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", err
	}

	return strings.TrimPrefix(release.TagName, "v"), nil
}

func compareVersions(current, latest string) bool {
	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	maxLen := len(currentParts)
	if len(latestParts) > maxLen {
		maxLen = len(latestParts)
	}

	for i := 0; i < maxLen; i++ {
		currentNum := 0
		latestNum := 0

		if i < len(currentParts) {
			fmt.Sscanf(currentParts[i], "%d", &currentNum)
		}
		if i < len(latestParts) {
			fmt.Sscanf(latestParts[i], "%d", &latestNum)
		}

		if latestNum > currentNum {
			return true
		} else if latestNum < currentNum {
			return false
		}
	}

	return false
}

func checkForUpdates() {
	latest, err := getLatestVersion()
	if err != nil {
		return
	}

	current := getCurrentVersion()
	if compareVersions(current, latest) {
		fmt.Printf("%s %s\n",
			colorize("🔄 Update available:", InfoColor),
			colorize(fmt.Sprintf("v%s → v%s", current, latest), SuccessColor))
		fmt.Printf("%s %s\n",
			colorize("📦 Run", InfoColor),
			colorize("echowave update", PrimaryColor)+colorize(" to update", InfoColor))
		fmt.Println()
	}
}

func getDownloadURL(release *GitHubRelease) string {
	var suffix string
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "amd64" {
			suffix = "macos-intel"
		} else {
			suffix = "macos-arm64"
		}
	case "linux":
		suffix = "linux-" + runtime.GOARCH
	case "windows":
		suffix = "windows-" + runtime.GOARCH + ".exe"
	default:
		return ""
	}

	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, suffix) {
			return asset.BrowserDownloadURL
		}
	}

	return ""
}

func performUpdate() {
	fmt.Printf("%s\n", colorize("🔄 Checking for updates...", InfoColor))

	latest, err := getLatestVersion()
	if err != nil {
		exitWithError(newError("Failed to check for updates", err))
	}

	current := getCurrentVersion()
	if !compareVersions(current, latest) {
		fmt.Printf("%s %s\n",
			colorize("✅ Already up to date:", SuccessColor),
			colorize("v"+current, PrimaryColor))
		return
	}

	fmt.Printf("%s %s\n",
		colorize("📦 Updating from", InfoColor),
		colorize(fmt.Sprintf("v%s to v%s", current, latest), PrimaryColor))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(GITHUB_API_URL)
	if err != nil {
		exitWithError(newError("Failed to fetch release information", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		exitWithError(newError("Failed to read release information", err))
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		exitWithError(newError("Failed to parse release information", err))
	}

	downloadURL := getDownloadURL(&release)
	if downloadURL == "" {
		exitWithError(newError("No compatible binary found for your platform", nil))
	}

	fmt.Printf("%s\n", colorize("⬇️ Downloading update...", InfoColor))

	resp, err = client.Get(downloadURL)
	if err != nil {
		exitWithError(newError("Failed to download update", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		exitWithError(newError(fmt.Sprintf("Failed to download update: %d", resp.StatusCode), nil))
	}

	execPath, err := os.Executable()
	if err != nil {
		exitWithError(newError("Failed to get executable path", err))
	}

	tempFile := execPath + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		exitWithError(newError("Failed to create temporary file", err))
	}

	_, err = io.Copy(file, resp.Body)
	file.Close()
	if err != nil {
		os.Remove(tempFile)
		exitWithError(newError("Failed to write update", err))
	}

	err = os.Chmod(tempFile, 0755)
	if err != nil {
		os.Remove(tempFile)
		exitWithError(newError("Failed to set executable permissions", err))
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", "timeout", "1", "&", "move", tempFile, execPath)
		cmd.Start()
	} else {
		err = os.Rename(tempFile, execPath)
		if err != nil {
			os.Remove(tempFile)
			exitWithError(newError("Failed to replace executable", err))
		}
	}

	fmt.Printf("%s %s\n",
		colorize("✅ Successfully updated to", SuccessColor),
		colorize("v"+latest, PrimaryColor))
	fmt.Printf("%s\n", colorize("🎉 Restart echowave to use the new version", InfoColor))
}

func showVersion() {
	fmt.Printf("%s %s\n",
		colorize("EchoWave", PrimaryColor),
		colorize("v"+getCurrentVersion(), SecondaryColor))
	fmt.Printf("%s %s\n",
		colorize("Part of the", MutedColor),
		betterLyrics()+colorize(" ecosystem", MutedColor))
	fmt.Printf("%s %s\n",
		colorize("Visit:", MutedColor),
		link("https://better-lyrics.boidu.dev"))
	os.Exit(0)
}

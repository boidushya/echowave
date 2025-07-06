package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

var VERSION = "dev"

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

func isNewerVersion(current, latest string) bool {
	if current == "dev" {
		return false
	}

	re := regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)`)

	currentMatch := re.FindStringSubmatch(current)
	latestMatch := re.FindStringSubmatch(latest)

	if len(currentMatch) < 4 || len(latestMatch) < 4 {
		return false
	}

	for i := 1; i <= 3; i++ {
		curr, _ := strconv.Atoi(currentMatch[i])
		lat, _ := strconv.Atoi(latestMatch[i])

		if lat > curr {
			return true
		}
		if lat < curr {
			return false
		}
	}

	return false
}

func getLatestRelease() (*GitHubRelease, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/repos/boidushya/echowave/releases/latest", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func findBinaryAsset(release *GitHubRelease) string {
	platform := runtime.GOOS
	arch := runtime.GOARCH

	expectedName := fmt.Sprintf("echowave-%s-%s", platform, arch)
	if platform == "windows" {
		expectedName += ".exe"
	}

	for _, asset := range release.Assets {
		if asset.Name == expectedName {
			return asset.URL
		}
	}

	return ""
}

func downloadBinary(url, path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func checkForUpdates() {
	if VERSION == "dev" {
		return
	}

	release, err := getLatestRelease()
	if err != nil {
		return
	}

	if isNewerVersion(VERSION, release.TagName) {
		fmt.Printf("%s %s\n",
			colorize("🔄 Update available:", InfoColor),
			colorize(fmt.Sprintf("v%s → %s", VERSION, release.TagName), SuccessColor))
		fmt.Printf("%s %s\n",
			colorize("📦 Run", InfoColor),
			colorize("echowave update", PrimaryColor)+colorize(" to update", InfoColor))
		fmt.Println()
	}
}

func performUpdate() {
	fmt.Printf("%s\n", colorize("🔄 Checking for updates...", InfoColor))

	if VERSION == "dev" {
		fmt.Printf("%s\n", colorize("⚠️  Development version - updates not available", WarningColor))
		return
	}

	release, err := getLatestRelease()
	if err != nil {
		exitWithError(newError("Failed to check for updates", err))
	}

	if !isNewerVersion(VERSION, release.TagName) {
		fmt.Printf("%s %s\n",
			colorize("✅ Already up to date:", SuccessColor),
			colorize("v"+VERSION, PrimaryColor))
		return
	}

	fmt.Printf("%s %s\n",
		colorize("📦 Updating from", InfoColor),
		colorize(fmt.Sprintf("v%s to %s", VERSION, release.TagName), PrimaryColor))

	downloadURL := findBinaryAsset(release)
	if downloadURL == "" {
		exitWithError(newError("No binary found for your platform", nil))
	}

	execPath, err := os.Executable()
	if err != nil {
		exitWithError(newError("Failed to get executable path", err))
	}

	tempPath := execPath + ".new"

	fmt.Printf("%s\n", colorize("⬇️  Downloading...", InfoColor))
	if err := downloadBinary(downloadURL, tempPath); err != nil {
		exitWithError(newError("Download failed", err))
	}

	if err := os.Chmod(tempPath, 0755); err != nil {
		os.Remove(tempPath)
		exitWithError(newError("Failed to set permissions", err))
	}

	fmt.Printf("%s\n", colorize("🔄 Installing...", InfoColor))
	if err := os.Rename(tempPath, execPath); err != nil {
		os.Remove(tempPath)
		exitWithError(newError("Failed to install update", err))
	}

	fmt.Printf("%s %s\n",
		colorize("✅ Updated to", SuccessColor),
		colorize(release.TagName, PrimaryColor))
}

func showVersion() {
	fmt.Printf("%s %s\n",
		colorize("EchoWave", PrimaryColor),
		colorize("v"+VERSION, SecondaryColor))
	fmt.Printf("%s %s\n",
		colorize("Part of the", MutedColor),
		betterLyrics()+colorize(" ecosystem", MutedColor))
	fmt.Printf("%s %s\n",
		colorize("Visit:", MutedColor),
		link("https://better-lyrics.boidu.dev"))
	os.Exit(0)
}

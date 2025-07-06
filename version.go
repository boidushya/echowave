package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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
		expectedName += ".zip"
	} else {
		expectedName += ".tar.gz"
	}

	for _, asset := range release.Assets {
		if asset.Name == expectedName {
			return asset.URL
		}
	}

	return ""
}

func downloadAndExtractBinary(url, targetPath string) error {
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

	tempDir, err := os.MkdirTemp("", "echowave-update")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if strings.HasSuffix(url, ".tar.gz") {
		return extractTarGz(resp.Body, tempDir, targetPath)
	} else if strings.HasSuffix(url, ".zip") {
		return extractZip(resp.Body, tempDir, targetPath)
	}

	return fmt.Errorf("unsupported archive format")
}

func extractTarGz(r io.Reader, tempDir, targetPath string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeReg {
			tempFile := filepath.Join(tempDir, filepath.Base(header.Name))
			file, err := os.Create(tempFile)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(file, tr); err != nil {
				return err
			}

			if err := os.Chmod(tempFile, 0755); err != nil {
				return err
			}

			return os.Rename(tempFile, targetPath)
		}
	}

	return fmt.Errorf("no binary found in archive")
}

func extractZip(r io.Reader, tempDir, targetPath string) error {
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	tempFile := filepath.Join(tempDir, "archive.zip")
	if err := os.WriteFile(tempFile, body, 0644); err != nil {
		return err
	}

	zr, err := zip.OpenReader(tempFile)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, file := range zr.File {
		if strings.HasSuffix(file.Name, ".exe") || !strings.Contains(file.Name, ".") {
			rc, err := file.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, rc); err != nil {
				return err
			}

			return os.Chmod(targetPath, 0755)
		}
	}

	return fmt.Errorf("no binary found in archive")
}

func isWritable(dir string) bool {
	testFile := filepath.Join(dir, ".echowave-write-test")
	file, err := os.Create(testFile)
	if err != nil {
		return false
	}
	file.Close()
	os.Remove(testFile)
	return true
}

func checkForUpdates() {
	if VERSION == "dev" {
		return
	}

	release, err := getLatestRelease()
	if err != nil {
		// Debug: show error if verbose mode or env var is set
		if os.Getenv("ECHOWAVE_DEBUG") == "1" {
			fmt.Printf("Debug: Failed to check for updates: %v\n", err)
		}
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
	if err := downloadAndExtractBinary(downloadURL, tempPath); err != nil {
		exitWithError(newError("Download failed", err))
	}

	fmt.Printf("%s\n", colorize("🔄 Installing...", InfoColor))
	
	// Check if we can write to the directory
	execDir := filepath.Dir(execPath)
	if !isWritable(execDir) {
		fmt.Printf("%s\n", colorize("⚠️  Root privileges required for installation", WarningColor))
		fmt.Printf("%s\n", colorize("Please run the following command:", InfoColor))
		fmt.Printf("%s\n", colorize(fmt.Sprintf("sudo cp %s %s", tempPath, execPath), PrimaryColor))
		fmt.Printf("%s\n", colorize("Then run: sudo chmod +x "+execPath, PrimaryColor))
		os.Exit(0)
	}
	
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

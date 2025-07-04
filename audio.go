// Package main implements EchoWave, a command-line tool for AI-powered audio transcription.
package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	downloadSpinnerDuration = 500 * time.Millisecond
	progressSpinnerDuration = 200 * time.Millisecond
)

var (
	ErrUnsupportedAudioFormat = errors.New("unsupported audio format")
	ErrAudioFileNotFound      = errors.New("audio file not found")
)

// isYouTubeURL validates if input string matches YouTube URL patterns.
// Supports both youtube.com and youtu.be domains with optional protocol and www prefix.
func isYouTubeURL(input string) bool {
	re := regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.be)/`)
	return re.MatchString(input)
}

// sanitizeYouTubeURL removes playlist parameters from YouTube URLs to ensure single video download.
// Strips list, index, and other playlist-related parameters while preserving the video ID.
func sanitizeYouTubeURL(input string) string {
	if !isYouTubeURL(input) {
		return input
	}

	parsedURL, err := url.Parse(input)
	if err != nil {
		return input
	}

	query := parsedURL.Query()

	// Remove playlist-related parameters
	query.Del("list")
	query.Del("index")
	query.Del("start_radio")
	query.Del("rv")

	// Keep only essential parameters (v for video ID, t for timestamp)
	sanitizedQuery := url.Values{}
	if v := query.Get("v"); v != "" {
		sanitizedQuery.Set("v", v)
	}
	if t := query.Get("t"); t != "" {
		sanitizedQuery.Set("t", t)
	}

	parsedURL.RawQuery = sanitizedQuery.Encode()
	return parsedURL.String()
}

// validateAudioFormat checks if the audio format is safe for command execution.
func validateAudioFormat(format string) bool {
	allowed := []string{"mp3", "wav", "m4a", "aac", "flac", "ogg"}
	format = strings.ToLower(strings.TrimSpace(format))
	for _, allowed := range allowed {
		if format == allowed {
			return true
		}
	}
	return false
}

// downloadYouTubeAudio extracts audio from YouTube URLs using yt-dlp.
// Creates temporary directory, downloads in specified format, and returns local file path.
// Verbose flag controls whether yt-dlp output is shown to user.
func downloadYouTubeAudio(url, audioFormat string, verbose bool) (string, error) {
	download("Downloading YouTube audio...")

	if !validateAudioFormat(audioFormat) {
		return "", newError("validate audio format", fmt.Errorf("%w: %s", ErrUnsupportedAudioFormat, audioFormat))
	}

	tmpDir, err := os.MkdirTemp("", "echowave-*")
	if err != nil {
		return "", newError("create temp directory", err)
	}

	outputPath := filepath.Join(tmpDir, "%(title)s.%(ext)s")
	cmd := exec.Command("yt-dlp", "-x", "--audio-format", audioFormat, "-o", outputPath, url)

	if !verbose {
		cmd.Stdout = nil
		cmd.Stderr = nil
	}

	done := make(chan bool)
	go func() {
		spinner("Downloading from YouTube...", downloadSpinnerDuration)
		for {
			select {
			case <-done:
				return
			default:
				spinner("Downloading from YouTube...", progressSpinnerDuration)
			}
		}
	}()

	err = cmd.Run()
	done <- true

	if err != nil {
		return "", newError("download YouTube audio", err)
	}

	pattern := fmt.Sprintf("*.%s", audioFormat)
	matches, err := filepath.Glob(filepath.Join(tmpDir, pattern))
	if err != nil || len(matches) == 0 {
		return "", newError("locate downloaded audio file", fmt.Errorf("%w: .%s", ErrAudioFileNotFound, audioFormat))
	}

	success("Audio download completed")
	return matches[0], nil
}

// processAudio determines whether input is YouTube URL or local file and handles accordingly.
// Returns audio file path and cleanup function for temporary files.
// Cleanup function removes temporary directories created for YouTube downloads.
func processAudio(input string, config *Config) (string, func()) {
	var audioPath string
	var cleanup func()

	if isYouTubeURL(input) {
		var err error
		sanitizedURL := sanitizeYouTubeURL(input)
		audioPath, err = downloadYouTubeAudio(sanitizedURL, config.AudioFormat, config.Verbose)
		if err != nil {
			exitWithError(err)
		}
		tempDir := filepath.Dir(audioPath)
		cleanup = func() {
			if err := os.RemoveAll(tempDir); err != nil {
				fmt.Printf("⚠️ Failed to cleanup temp files: %v\n", err)
			}
		}
	} else {
		audioPath = input
		cleanup = func() {}
	}

	return audioPath, cleanup
}

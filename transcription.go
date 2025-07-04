// Package main implements EchoWave, a command-line tool for AI-powered audio transcription.
// It converts audio files and YouTube URLs into synchronized lyrics using OpenAI Whisper.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	secondsPerMinute = 60
)

var (
	ErrUnsupportedWhisperModel = errors.New("unsupported whisper model")
	ErrNoSegmentsFound         = errors.New("no segments found in transcription")
)

// Word represents a single word with timing and confidence information.
// Used for word-level accuracy visualization in heatmap mode.
type Word struct {
	Word        string  `json:"word"`
	Start       float64 `json:"start"`
	End         float64 `json:"end"`
	Probability float64 `json:"probability"`
}

// Segment represents a single transcribed text segment with timing information.
// Used for parsing Whisper JSON output and generating LRC timestamps.
type Segment struct {
	Start       float64 `json:"start"`
	Text        string  `json:"text"`
	AvgLogprob  float64 `json:"avg_logprob"`
	Confidence  float64 `json:"confidence"`
	Words       []Word  `json:"words"`
}

// WhisperOutput represents the complete JSON response from OpenAI Whisper transcription.
// Contains an array of text segments with precise timing for lyrics generation.
type WhisperOutput struct {
	Segments []Segment `json:"segments"`
}

// secondsToLRCTimestamp converts floating-point seconds to LRC synchronized lyric format [MM:SS.XX].
// Used for creating timestamps compatible with media players that support LRC files.
func secondsToLRCTimestamp(seconds float64) string {
	minutes := int(seconds) / secondsPerMinute
	sec := seconds - float64(minutes*secondsPerMinute)
	return fmt.Sprintf("[%02d:%05.2f]", minutes, sec)
}

// getConfidenceColor returns the appropriate color based on confidence level.
// High confidence (>0.8): Green, Medium confidence (0.5-0.8): Yellow, Low confidence (<0.5): Red
func getConfidenceColor(confidence float64) string {
	if confidence >= 0.8 {
		return BrightGreen
	} else if confidence >= 0.5 {
		return BrightYellow
	}
	return BrightRed
}

// displayHeatmap shows a color-coded visualization of transcription accuracy.
// Words are colored based on their confidence scores for easy identification of uncertain transcription.
func displayHeatmap(jsonPath string) error {
	header("Transcription Accuracy Heatmap")
	
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return newError("read JSON file for heatmap", err)
	}

	var output WhisperOutput
	if err := json.Unmarshal(data, &output); err != nil {
		return newError("parse JSON for heatmap", err)
	}

	if len(output.Segments) == 0 {
		return newError("display heatmap", ErrNoSegmentsFound)
	}

	info("Legend: " + colorize("High confidence (>0.8)", BrightGreen) + " | " + 
		colorize("Medium confidence (0.5-0.8)", BrightYellow) + " | " + 
		colorize("Low confidence (<0.5)", BrightRed))
	fmt.Println()

	for _, segment := range output.Segments {
		fmt.Printf("%s ", colorize(secondsToLRCTimestamp(segment.Start), MutedColor))
		
		if len(segment.Words) > 0 {
			for _, word := range segment.Words {
				color := getConfidenceColor(word.Probability)
				fmt.Printf("%s ", colorize(word.Word, color))
			}
		} else {
			segmentConfidence := segment.Confidence
			if segmentConfidence == 0 {
				segmentConfidence = 1.0 + segment.AvgLogprob
				if segmentConfidence < 0 {
					segmentConfidence = 0
				}
			}
			color := getConfidenceColor(segmentConfidence)
			fmt.Printf("%s ", colorize(strings.TrimSpace(segment.Text), color))
		}
		fmt.Println()
	}

	fmt.Println()
	success("Heatmap display completed")
	return nil
}

// validateWhisperModel checks if the model is a valid Whisper model.
func validateWhisperModel(model string) bool {
	allowed := []string{"tiny", "base", "small", "medium", "large", "large-v1", "large-v2", "large-v3"}
	for _, allowed := range allowed {
		if model == allowed {
			return true
		}
	}
	return false
}

// runWhisper executes OpenAI Whisper AI transcription engine with audio file and model configuration.
// Outputs JSON transcription with word-level timestamps to specified directory.
// Stdout/stderr are inherited to show real-time transcription progress.
func runWhisper(audioPath, model, language, outputDir string) error {
	processing("Running Whisper transcription...")

	if !validateWhisperModel(model) {
		return newError("validate whisper model", fmt.Errorf("%w: %s", ErrUnsupportedWhisperModel, model))
	}

	step("Model: " + model + ", Language: " + language)

	cmd := exec.Command("whisper", audioPath, "--model", model, "--language", language,
		"--output_format", "json", "--word_timestamps", "True", "--temperature", "0", "--output_dir", outputDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return newError("run Whisper transcription", err)
	}

	success("Whisper transcription completed")
	return nil
}

// convertJSONToLRC parses Whisper's JSON transcription output and generates LRC synchronized lyrics file.
// Each segment's start timestamp is converted to LRC format with corresponding text.
// Validates JSON structure and ensures segments exist before processing.
func convertJSONToLRC(jsonPath, lrcPath string) error {
	step("Converting transcription to LRC format...")

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return newError("read JSON file", err)
	}

	var output WhisperOutput
	if err := json.Unmarshal(data, &output); err != nil {
		return newError("parse JSON", err)
	}

	if len(output.Segments) == 0 {
		return newError("process transcription", ErrNoSegmentsFound)
	}

	lrcFile, err := os.Create(lrcPath)
	if err != nil {
		return newError("create LRC file", err)
	}
	defer func() {
		if err := lrcFile.Close(); err != nil {
			fmt.Printf("Warning: failed to close LRC file: %v\n", err)
		}
	}()

	for _, segment := range output.Segments {
		line := fmt.Sprintf("%s %s\n", secondsToLRCTimestamp(segment.Start), strings.TrimSpace(segment.Text))
		if _, err := lrcFile.WriteString(line); err != nil {
			return newError("write LRC content", err)
		}
	}

	file("LRC file created: " + lrcPath)
	return nil
}

// generateTranscription manages the complete audio-to-lyrics pipeline using Whisper AI.
// Creates output directory, runs transcription, handles file naming, and generates both JSON and LRC formats.
// Automatically resolves output file paths and manages temporary file cleanup.
func generateTranscription(audioPath string, config *Config) {
	step("Setting up output directory...")
	if err := os.MkdirAll(config.OutputDir, 0o750); err != nil {
		exitWithError(newError("create output directory", err))
	}

	var base string
	if config.Output != "" {
		base = filepath.Join(config.OutputDir, config.Output)
	} else {
		audioBaseName := strings.TrimSuffix(filepath.Base(audioPath), filepath.Ext(audioPath))
		base = filepath.Join(config.OutputDir, audioBaseName)
	}

	jsonPath := base + ".json"
	lrcPath := base + ".lrc"

	if err := runWhisper(audioPath, config.Model, config.Language, config.OutputDir); err != nil {
		exitWithError(err)
	}

	actualJSONPath := jsonPath
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		audioBaseName := strings.TrimSuffix(filepath.Base(audioPath), filepath.Ext(audioPath))
		whisperJSONPath := filepath.Join(config.OutputDir, audioBaseName+".json")
		if _, err := os.Stat(whisperJSONPath); err == nil {
			actualJSONPath = whisperJSONPath
		}
	}

	if err := convertJSONToLRC(actualJSONPath, lrcPath); err != nil {
		exitWithError(err)
	}

	if config.Heatmap {
		fmt.Println()
		if err := displayHeatmap(actualJSONPath); err != nil {
			warning("Failed to display heatmap: " + err.Error())
		}
	}

	fmt.Println()
	success("Transcription completed successfully!")
	info("Files saved in: " + config.OutputDir)
}

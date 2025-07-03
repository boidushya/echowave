package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Segment struct {
	Start float64 `json:"start"`
	Text  string  `json:"text"`
}

type WhisperOutput struct {
	Segments []Segment `json:"segments"`
}

type Config struct {
	Model       string
	Language    string
	AudioFormat string
	OutputDir   string
	Output      string
}

func secondsToLRCTimestamp(seconds float64) string {
	min := int(seconds) / 60
	sec := seconds - float64(min*60)
	return fmt.Sprintf("[%02d:%05.2f]", min, sec)
}

func isYouTubeURL(input string) bool {
	re := regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.be)/`)
	return re.MatchString(input)
}

func downloadYouTubeAudio(url, audioFormat string) (string, error) {
	download("Downloading YouTube audio...")

	tmpDir, err := os.MkdirTemp("", "echowave-*")
	if err != nil {
		return "", newError("create temp directory", err)
	}

	outputPath := filepath.Join(tmpDir, "%(title)s.%(ext)s")
	cmd := exec.Command("yt-dlp", "-x", "--audio-format", audioFormat, "-o", outputPath, url)
	
	// Hide yt-dlp output for cleaner experience
	cmd.Stdout = nil
	cmd.Stderr = nil

	// Show spinner while downloading
	done := make(chan bool)
	go func() {
		spinner("Downloading from YouTube...", 500*time.Millisecond)
		for {
			select {
			case <-done:
				return
			default:
				spinner("Downloading from YouTube...", 200*time.Millisecond)
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
		return "", newError("locate downloaded audio file", fmt.Errorf("could not find .%s file", audioFormat))
	}
	
	success("Audio download completed")
	return matches[0], nil
}

func runWhisper(audioPath, model, language, outputDir string) error {
	processing("Running Whisper transcription...")
	step("Model: " + model + ", Language: " + language)
	
	cmd := exec.Command("whisper", audioPath, "--model", model, "--language", language, "--output_format", "json", "--word_timestamps", "True", "--temperature", "0", "--output_dir", outputDir)
	
	// Hide whisper output for cleaner experience
	cmd.Stdout = nil
	cmd.Stderr = nil
	
	// Show spinner while processing
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				spinner("Processing with Whisper AI...", 300*time.Millisecond)
			}
		}
	}()
	
	err := cmd.Run()
	done <- true
	
	if err != nil {
		return newError("run Whisper transcription", err)
	}
	
	success("Whisper transcription completed")
	return nil
}

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
		return newError("process transcription", fmt.Errorf("no segments found in JSON"))
	}

	lrcFile, err := os.Create(lrcPath)
	if err != nil {
		return newError("create LRC file", err)
	}
	defer lrcFile.Close()

	for _, segment := range output.Segments {
		line := fmt.Sprintf("%s %s\n", secondsToLRCTimestamp(segment.Start), strings.TrimSpace(segment.Text))
		if _, err := lrcFile.WriteString(line); err != nil {
			return newError("write LRC content", err)
		}
	}

	file("LRC file created: " + lrcPath)
	return nil
}

func parseFlags() *Config {
	var (
		model       = flag.String("model", "large-v3", "Whisper model to use")
		language    = flag.String("language", "en", "Language for transcription")
		audioFormat = flag.String("audio-format", "mp3", "Audio format for download")
		outputDir   = flag.String("output-dir", "dist", "Output directory for generated files")
		output      = flag.String("output", "", "Output file path (without extension)")
		help        = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		fmt.Print(logo())
		fmt.Println()
		
		header("EchoWave - Audio Transcription Tool")
		subheader("Part of the " + betterLyrics() + " ecosystem")
		fmt.Println()
		
		info("Transform audio into lyrics with AI-powered transcription")
		fmt.Printf("%s%s\n", prefix(), colorize("🌐 Visit: ", InfoColor) + link("https://better-lyrics.boidu.dev"))
		fmt.Println()
		
		fmt.Print(box("Usage", "echowave [OPTIONS] <YouTube URL or path/to/audio>"))
		fmt.Println()
		
		header("Options")
		fmt.Printf("%s%s\n", prefix(), colorize("  -model string", PrimaryColor))
		fmt.Printf("%s%s\n", prefix(), colorize("        Whisper model to use (default \"large-v3\")", MutedColor))
		fmt.Printf("%s%s\n", prefix(), colorize("  -language string", PrimaryColor))
		fmt.Printf("%s%s\n", prefix(), colorize("        Language for transcription (default \"en\")", MutedColor))
		fmt.Printf("%s%s\n", prefix(), colorize("  -audio-format string", PrimaryColor))
		fmt.Printf("%s%s\n", prefix(), colorize("        Audio format for download (default \"mp3\")", MutedColor))
		fmt.Printf("%s%s\n", prefix(), colorize("  -output-dir string", PrimaryColor))
		fmt.Printf("%s%s\n", prefix(), colorize("        Output directory for generated files (default \"dist\")", MutedColor))
		fmt.Printf("%s%s\n", prefix(), colorize("  -output string", PrimaryColor))
		fmt.Printf("%s%s\n", prefix(), colorize("        Output file path (without extension)", MutedColor))
		fmt.Printf("%s%s\n", prefix(), colorize("  -help", PrimaryColor))
		fmt.Printf("%s%s\n", prefix(), colorize("        Show this help message", MutedColor))
		fmt.Println()
		
		header("Examples")
		fmt.Printf("%s%s\n", prefix(), colorize("# Transcribe YouTube video", SecondaryColor))
		fmt.Printf("%s%s\n", prefix(), colorize("echowave https://youtube.com/watch?v=xyz", White))
		fmt.Println()
		fmt.Printf("%s%s\n", prefix(), colorize("# Transcribe local audio file", SecondaryColor))
		fmt.Printf("%s%s\n", prefix(), colorize("echowave audio.mp3", White))
		fmt.Println()
		fmt.Printf("%s%s\n", prefix(), colorize("# Custom model and language", SecondaryColor))
		fmt.Printf("%s%s\n", prefix(), colorize("echowave -model=medium -language=es -output=transcript audio.mp3", White))
		fmt.Println()
		fmt.Printf("%s%s\n", prefix(), colorize("# Custom output directory", SecondaryColor))
		fmt.Printf("%s%s\n", prefix(), colorize("echowave -output-dir=transcripts https://youtube.com/watch?v=xyz", White))
		fmt.Println()
		
		fmt.Printf("%s%s\n", prefix(), colorize("Made with ❤️ by the ", MutedColor) + betterLyrics() + colorize(" team", MutedColor))
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		errorMsg("Please provide a YouTube URL or audio file path")
		info("Use -help for usage information")
		os.Exit(1)
	}

	return &Config{
		Model:       *model,
		Language:    *language,
		AudioFormat: *audioFormat,
		OutputDir:   *outputDir,
		Output:      *output,
	}
}

func processAudio(input string, config *Config) (string, func()) {
	var audioPath string
	var cleanup func()

	if isYouTubeURL(input) {
		var err error
		audioPath, err = downloadYouTubeAudio(input, config.AudioFormat)
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

func generateTranscription(audioPath string, config *Config) {
	step("Setting up output directory...")
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
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

	fmt.Println()
	success("Transcription completed successfully!")
	info("Files saved in: " + config.OutputDir)
}

func main() {
	config := parseFlags()
	
	if !checkAllDependencies() {
		os.Exit(1)
	}
	
	input := flag.Arg(0)

	audioPath, cleanup := processAudio(input, config)
	defer cleanup()

	generateTranscription(audioPath, config)
}

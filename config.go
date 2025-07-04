package main

import (
	"flag"
	"fmt"
	"os"
)

// Config holds all command-line configuration options for EchoWave transcription.
// Contains Whisper model settings, audio processing options, and output preferences.
type Config struct {
	Model       string
	Language    string
	AudioFormat string
	OutputDir   string
	Output      string
	Verbose     bool
	Heatmap     bool
}

// showHelp displays the complete help documentation for EchoWave including usage examples,
// command-line options, and installation links. This function presents a formatted help screen
// with colored output and ASCII art branding, then exits the program with status code 0.
// The help content includes practical examples for common use cases like YouTube transcription,
// local file processing, and custom model configuration.
func showHelp() {
	fmt.Print(logo())
	fmt.Println()

	fmt.Printf("%s\n", colorize(bold("EchoWave - Audio Transcription Tool"), PrimaryColor))
	fmt.Printf("%s\n", colorize("Part of the "+betterLyrics()+" ecosystem", SecondaryColor))
	fmt.Println()

	description := "Transform audio into lyrics with AI-powered transcription"
	fmt.Printf("%s %s\n", colorize("ℹ️", InfoColor), colorize(description, InfoColor))
	fmt.Printf("%s %s\n", colorize("🌐", InfoColor), colorize("Visit: ", InfoColor)+link("https://better-lyrics.boidu.dev"))
	fmt.Println()

	fmt.Print(box("Usage", "echowave [OPTIONS] <YouTube URL or path/to/audio>"))
	fmt.Println()

	fmt.Printf("%s\n", colorize(bold("Options"), PrimaryColor))
	fmt.Printf("%s\n", colorize("  -model string", PrimaryColor))
	fmt.Printf("%s\n", colorize("        Whisper model to use (default \"large-v3\")", MutedColor))
	fmt.Printf("%s\n", colorize("  -language string", PrimaryColor))
	fmt.Printf("%s\n", colorize("        Language for transcription (default \"en\")", MutedColor))
	fmt.Printf("%s\n", colorize("  -audio-format string", PrimaryColor))
	fmt.Printf("%s\n", colorize("        Audio format for download (default \"mp3\")", MutedColor))
	fmt.Printf("%s\n", colorize("  -output-dir string", PrimaryColor))
	fmt.Printf("%s\n", colorize("        Output directory for generated files (default \".\")", MutedColor))
	fmt.Printf("%s\n", colorize("  -output string", PrimaryColor))
	fmt.Printf("%s\n", colorize("        Output file path (without extension)", MutedColor))
	fmt.Printf("%s\n", colorize("  -verbose", PrimaryColor))
	fmt.Printf("%s\n", colorize("        Show detailed output from tools", MutedColor))
	fmt.Printf("%s\n", colorize("  -heatmap", PrimaryColor))
	fmt.Printf("%s\n", colorize("        Show transcription accuracy heatmap (default true)", MutedColor))
	fmt.Printf("%s\n", colorize("  -help", PrimaryColor))
	fmt.Printf("%s\n", colorize("        Show this help message", MutedColor))
	fmt.Println()

	fmt.Printf("%s\n", colorize(bold("Examples"), PrimaryColor))
	fmt.Printf("%s\n", colorize("# Transcribe YouTube video", SecondaryColor))
	fmt.Printf("%s\n", colorize("echowave https://youtube.com/watch?v=xyz", White))
	fmt.Println()
	fmt.Printf("%s\n", colorize("# Transcribe local audio file", SecondaryColor))
	fmt.Printf("%s\n", colorize("echowave audio.mp3", White))
	fmt.Println()
	fmt.Printf("%s\n", colorize("# Custom model and language", SecondaryColor))
	fmt.Printf("%s\n", colorize("echowave -model=medium -language=es -output=transcript audio.mp3", White))
	fmt.Println()
	fmt.Printf("%s\n", colorize("# Custom output directory", SecondaryColor))
	fmt.Printf("%s\n", colorize("echowave -output-dir=transcripts https://youtube.com/watch?v=xyz", White))
	fmt.Println()
	fmt.Printf("%s\n", colorize("# Verbose output", SecondaryColor))
	fmt.Printf("%s\n", colorize("echowave -verbose audio.mp3", White))
	fmt.Println()
	fmt.Printf("%s\n", colorize("# Disable accuracy heatmap", SecondaryColor))
	fmt.Printf("%s\n", colorize("echowave -heatmap=false audio.mp3", White))
	fmt.Println()

	fmt.Printf("%s\n", colorize("Made with ❤️ by the ", MutedColor)+betterLyrics()+colorize(" team", MutedColor))
	os.Exit(0)
}

// parseFlags processes command-line arguments and returns a populated Config struct.
// It defines and parses all supported flags including model selection, language settings,
// audio format preferences, output directory, and verbose mode. If the help flag is set
// or no arguments are provided, it automatically displays help and exits. The function
// validates that at least one positional argument (audio source) is provided before
// returning the configuration object.
func parseFlags() *Config {
	var (
		model       = flag.String("model", "medium", "Whisper model to use")
		language    = flag.String("language", "en", "Language for transcription")
		audioFormat = flag.String("audio-format", "mp3", "Audio format for download")
		outputDir   = flag.String("output-dir", ".", "Output directory for generated files")
		output      = flag.String("output", "", "Output file path (without extension)")
		verbose     = flag.Bool("verbose", false, "Show detailed output from tools")
		heatmap     = flag.Bool("heatmap", true, "Show transcription accuracy heatmap")
		help        = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
	}

	if flag.NArg() < 1 {
		showHelp()
	}

	return &Config{
		Model:       *model,
		Language:    *language,
		AudioFormat: *audioFormat,
		OutputDir:   *outputDir,
		Output:      *output,
		Verbose:     *verbose,
		Heatmap:     *heatmap,
	}
}

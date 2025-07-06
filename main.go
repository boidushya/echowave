package main

import (
	"flag"
	"os"
)

// main orchestrates the complete EchoWave audio transcription workflow from start to finish.
// It parses command-line flags, validates all required dependencies are installed, processes
// the input audio source (either downloading from YouTube or using a local file), and generates
// the final transcription output. The function handles cleanup of temporary files through defer
// statements and exits with appropriate error codes if any step fails. This is the primary
// entry point that coordinates all other application components.
func main() {
	config := parseFlags()

	checkForUpdates()

	if !checkAllDependencies() {
		os.Exit(1)
	}

	input := flag.Arg(0)

	audioPath, cleanup := processAudio(input, config)
	defer cleanup()

	generateTranscription(audioPath, config)
}

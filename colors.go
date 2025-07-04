package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	spinnerSleepDuration = 100 * time.Millisecond
	clearLinePadding     = 20
	boxPadding           = 4
	boxBorderPadding     = 2
	boxContentPadding    = 3
)

// ANSI color codes.
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"

	// Colors.
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Bright colors.
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	// Background colors.
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

// Theme colors.
var (
	PrimaryColor   = BrightBlue
	SecondaryColor = BrightCyan
	SuccessColor   = BrightGreen
	WarningColor   = BrightYellow
	ErrorColor     = BrightRed
	InfoColor      = BrightMagenta
	MutedColor     = BrightBlack
	BrandColor     = BrightRed
)

// gradientColors defines 256-color ANSI escape sequences for logo gradient effect.
// Progresses from bright blue through cyan tones to light blue for visual appeal.
var gradientColors = []string{
	"\033[38;5;39m",
	"\033[38;5;45m",
	"\033[38;5;51m",
	"\033[38;5;87m",
	"\033[38;5;123m",
	"\033[38;5;159m",
}

func colorize(text, color string) string {
	return color + text + Reset
}

func bold(text string) string {
	return Bold + text + Reset
}

func underline(text string) string {
	return Underline + text + Reset
}

func prefix() string {
	return colorize("[", SecondaryColor) +
		colorize("echo", PrimaryColor) +
		colorize("wave", SuccessColor) +
		colorize("]", SecondaryColor) + " "
}

func logo() string {
	logoLines := []string{
		"  ███████╗ ██████╗██╗  ██╗ ██████╗ ██╗    ██╗ █████╗ ██╗   ██╗███████╗",
		"  ██╔════╝██╔════╝██║  ██║██╔═══██╗██║    ██║██╔══██╗██║   ██║██╔════╝",
		"  █████╗  ██║     ███████║██║   ██║██║ █╗ ██║███████║██║   ██║█████╗  ",
		"  ██╔══╝  ██║     ██╔══██║██║   ██║██║███╗██║██╔══██║╚██╗ ██╔╝██╔══╝  ",
		"  ███████╗╚██████╗██║  ██║╚██████╔╝╚███╔███╔╝██║  ██║ ╚████╔╝ ███████╗",
		"  ╚══════╝ ╚═════╝╚═╝  ╚═╝ ╚═════╝  ╚══╝╚══╝ ╚═╝  ╚═╝  ╚═══╝  ╚══════╝",
	}

	var result strings.Builder
	for i, line := range logoLines {
		colorIndex := i % len(gradientColors)
		result.WriteString(gradientColors[colorIndex] + line + Reset + "\n")
	}

	return result.String()
}

func spinner(message string, duration time.Duration) {
	spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	fmt.Print(prefix() + colorize(message, InfoColor) + " ")

	start := time.Now()
	i := 0
	for time.Since(start) < duration {
		fmt.Printf("\r%s%s %s", prefix(), colorize(message, InfoColor), colorize(spinChars[i%len(spinChars)], PrimaryColor))
		time.Sleep(spinnerSleepDuration)
		i++
	}
	fmt.Print("\r" + strings.Repeat(" ", len(message)+clearLinePadding) + "\r")
}

func success(message string) {
	fmt.Printf("%s%s %s\n", prefix(), colorize("✅", SuccessColor), colorize(message, SuccessColor))
}

func warning(message string) {
	fmt.Printf("%s%s %s\n", prefix(), colorize("⚠️", WarningColor), colorize(message, WarningColor))
}

func errorMsg(message string) {
	fmt.Printf("%s%s %s\n", prefix(), colorize("❌", ErrorColor), colorize(message, ErrorColor))
}

func info(message string) {
	fmt.Printf("%s%s %s\n", prefix(), colorize("ℹ️", InfoColor), colorize(message, InfoColor))
}

func step(message string) {
	fmt.Printf("%s%s %s\n", prefix(), colorize("🔄", PrimaryColor), colorize(message, PrimaryColor))
}

func download(message string) {
	fmt.Printf("%s%s %s\n", prefix(), colorize("📥", SecondaryColor), colorize(message, SecondaryColor))
}

func processing(message string) {
	fmt.Printf("%s%s %s\n", prefix(), colorize("🧠", InfoColor), colorize(message, InfoColor))
}

func file(message string) {
	fmt.Printf("%s%s %s\n", prefix(), colorize("📄", MutedColor), colorize(message, White))
}

func header(message string) {
	fmt.Printf("\n%s%s\n", prefix(), colorize(bold(message), PrimaryColor))
}

func subheader(message string) {
	fmt.Printf("%s%s\n", prefix(), colorize(message, SecondaryColor))
}

func betterLyrics() string {
	return colorize("better-lyrics", BrandColor)
}

func link(url string) string {
	return underline(colorize(url, BrightCyan))
}

// Box drawing characters for fancy borders.
func box(title, content string) string {
	lines := strings.Split(content, "\n")
	maxWidth := 0

	// Find the maximum width
	if len(title) > maxWidth {
		maxWidth = len(title)
	}
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Add padding
	width := maxWidth + boxPadding

	var result strings.Builder

	// Top border
	topBorder := colorize("╭", PrimaryColor) + colorize(strings.Repeat("─", width-2), PrimaryColor) + colorize("╮", PrimaryColor) + "\n"
	result.WriteString(topBorder)

	// Title
	if title != "" {
		padding := (width - len(title) - boxBorderPadding) / boxBorderPadding
		titleLine := colorize("│", PrimaryColor) + strings.Repeat(" ", padding) + colorize(bold(title), PrimaryColor) + strings.Repeat(" ", width-len(title)-padding-2) + colorize("│", PrimaryColor) + "\n"
		result.WriteString(titleLine)
		midBorder := colorize("├", PrimaryColor) + colorize(strings.Repeat("─", width-2), PrimaryColor) + colorize("┤", PrimaryColor) + "\n"
		result.WriteString(midBorder)
	}

	// Content
	for _, line := range lines {
		if line == "" {
			result.WriteString(colorize("│", PrimaryColor) + strings.Repeat(" ", width-boxBorderPadding) + colorize("│", PrimaryColor) + "\n")
		} else {
			padding := width - len(line) - boxContentPadding
			contentLine := colorize("│", PrimaryColor) + " " + line + strings.Repeat(" ", padding) + colorize("│", PrimaryColor) + "\n"
			result.WriteString(contentLine)
		}
	}

	// Bottom border
	bottomBorder := colorize("╰", PrimaryColor) + colorize(strings.Repeat("─", width-2), PrimaryColor) + colorize("╯", PrimaryColor) + "\n"
	result.WriteString(bottomBorder)

	return result.String()
}

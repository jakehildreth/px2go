package render

import (
	"fmt"
	"os"
	"runtime"
)

// TrueColorFg returns the ANSI escape sequence to set foreground to the given RGB color.
func TrueColorFg(r, g, b uint8) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

// TrueColorBg returns the ANSI escape sequence to set background to the given RGB color.
func TrueColorBg(r, g, b uint8) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
}

// AnsiReset returns the ANSI escape sequence to reset all text formatting.
func AnsiReset() string {
	return "\x1b[0m"
}

// DetectColorMode returns the appropriate color mode for the current environment.
// Priority: NO_COLOR → "none"; COLORTERM=truecolor|24bit → "truecolor";
// Windows (VT enabled) → "truecolor"; default → "truecolor".
func DetectColorMode() string {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return "none"
	}
	colorterm := os.Getenv("COLORTERM")
	if colorterm == "truecolor" || colorterm == "24bit" {
		return "truecolor"
	}
	if runtime.GOOS == "windows" {
		if err := enableWindowsVT(); err != nil {
			return "none"
		}
	}
	return "truecolor"
}

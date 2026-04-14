package render_test

import (
	"os"
	"testing"

	"github.com/jakehildreth/px2go/render"
)

func TestTrueColorFg(t *testing.T) {
	got := render.TrueColorFg(255, 128, 0)
	want := "\x1b[38;2;255;128;0m"
	if got != want {
		t.Errorf("TrueColorFg: got %q, want %q", got, want)
	}
}

func TestTrueColorBg(t *testing.T) {
	got := render.TrueColorBg(0, 64, 128)
	want := "\x1b[48;2;0;64;128m"
	if got != want {
		t.Errorf("TrueColorBg: got %q, want %q", got, want)
	}
}

func TestAnsiReset(t *testing.T) {
	got := render.AnsiReset()
	want := "\x1b[0m"
	if got != want {
		t.Errorf("AnsiReset: got %q, want %q", got, want)
	}
}

func TestDetectColorMode_NoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	got := render.DetectColorMode()
	if got != "none" {
		t.Errorf("with NO_COLOR set: got %q, want %q", got, "none")
	}
}

func TestDetectColorMode_Truecolor(t *testing.T) {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		t.Skip("NO_COLOR is set in environment")
	}
	t.Setenv("COLORTERM", "truecolor")
	got := render.DetectColorMode()
	if got != "truecolor" {
		t.Errorf("with COLORTERM=truecolor: got %q, want %q", got, "truecolor")
	}
}

func TestDetectColorMode_24bit(t *testing.T) {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		t.Skip("NO_COLOR is set in environment")
	}
	t.Setenv("COLORTERM", "24bit")
	got := render.DetectColorMode()
	if got != "truecolor" {
		t.Errorf("with COLORTERM=24bit: got %q, want %q", got, "truecolor")
	}
}

package render_test

import (
	"strings"
	"testing"

	"github.com/jakehildreth/px2go/render"
)

// oddHeight: 1x1 green — Y=-1 trick makes pixel land in bottom half → ▄ with FG=green
func TestRender_OddHeight(t *testing.T) {
	pixels := [][]byte{
		{0, 255, 0, 255}, // green opaque
	}

	var buf strings.Builder
	if err := render.Render(1, 1, pixels, "truecolor", &buf); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "▄") {
		t.Errorf("expected lower half block '▄' for 1px-height image; got: %q", out)
	}
	if !strings.Contains(out, "\x1b[38;2;0;255;0m") {
		t.Errorf("expected green FG sequence; got: %q", out)
	}
}

// evenHeight: 2x2, row 0 = all red, row 1 = all blue → BG=red, FG=blue, char=▄
func TestRender_EvenHeight_BothOpaque(t *testing.T) {
	pixels := [][]byte{
		{255, 0, 0, 255}, // (0,0) red
		{255, 0, 0, 255}, // (1,0) red
		{0, 0, 255, 255}, // (0,1) blue
		{0, 0, 255, 255}, // (1,1) blue
	}

	var buf strings.Builder
	if err := render.Render(2, 2, pixels, "truecolor", &buf); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "▄") {
		t.Errorf("expected lower half block '▄'; got: %q", out)
	}
	if !strings.Contains(out, "\x1b[48;2;255;0;0m") {
		t.Errorf("expected red BG sequence; got: %q", out)
	}
	if !strings.Contains(out, "\x1b[38;2;0;0;255m") {
		t.Errorf("expected blue FG sequence; got: %q", out)
	}
	if !strings.Contains(out, "\x1b[0m\x1b[K") {
		t.Errorf("expected reset+erase-to-EOL at end of row; got: %q", out)
	}
}

// top opaque red, bottom transparent → ▀ with FG=red
func TestRender_TopOpaqueBottomTransparent(t *testing.T) {
	pixels := [][]byte{
		{255, 0, 0, 255}, // top: red opaque
		{0, 0, 0, 0},     // bottom: transparent
	}

	var buf strings.Builder
	if err := render.Render(1, 2, pixels, "truecolor", &buf); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "▀") {
		t.Errorf("expected upper half block '▀'; got: %q", out)
	}
	if !strings.Contains(out, "\x1b[38;2;255;0;0m") {
		t.Errorf("expected red FG sequence; got: %q", out)
	}
}

// top transparent, bottom opaque blue → ▄ with FG=blue
func TestRender_TopTransparentBottomOpaque(t *testing.T) {
	pixels := [][]byte{
		{0, 0, 0, 0},     // top: transparent
		{0, 0, 255, 255}, // bottom: blue opaque
	}

	var buf strings.Builder
	if err := render.Render(1, 2, pixels, "truecolor", &buf); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "▄") {
		t.Errorf("expected lower half block '▄'; got: %q", out)
	}
	if !strings.Contains(out, "\x1b[38;2;0;0;255m") {
		t.Errorf("expected blue FG sequence; got: %q", out)
	}
}

// both transparent → space, no block chars
func TestRender_BothTransparent(t *testing.T) {
	pixels := [][]byte{
		{0, 0, 0, 0}, // transparent
		{0, 0, 0, 0}, // transparent
	}

	var buf strings.Builder
	if err := render.Render(1, 2, pixels, "truecolor", &buf); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()

	if strings.Contains(out, "▄") || strings.Contains(out, "▀") {
		t.Errorf("expected no block chars for fully transparent pixels; got: %q", out)
	}
	if !strings.Contains(out, " ") {
		t.Errorf("expected space for transparent cell; got: %q", out)
	}
}

// colorMode=none: no ANSI sequences, block chars still present
func TestRender_ColorModeNone(t *testing.T) {
	pixels := [][]byte{
		{255, 0, 0, 255}, // red opaque
		{0, 0, 255, 255}, // blue opaque
	}

	var buf strings.Builder
	if err := render.Render(1, 2, pixels, "none", &buf); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()

	if strings.Contains(out, "\x1b[") {
		t.Errorf("expected no ANSI sequences with colorMode=none; got: %q", out)
	}
	if !strings.Contains(out, "▄") {
		t.Errorf("expected block char even with colorMode=none; got: %q", out)
	}
}

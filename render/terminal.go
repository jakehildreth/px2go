package render

import (
	"fmt"
	"io"
	"strings"
)

const (
	lowerHalf = '▄' // U+2584: lower half block
	upperHalf = '▀' // U+2580: upper half block
)

// Render writes pixel art to out using terminal block characters.
// Two pixel rows are combined into one terminal row via ▄/▀ half-block characters.
// For odd-height images, rendering starts at virtual row Y=-1 (transparent) so the
// first actual pixel row appears as the bottom half of the first terminal row.
// colorMode must be "truecolor" or "none".
func Render(w, h uint32, pixels [][]byte, colorMode string, out io.Writer) error {
	useColor := colorMode != "none"

	startY := 0
	endY := int(h)
	if h%2 == 1 {
		startY = -1
		endY = int(h) - 1
	}

	for y := startY; y < endY; y += 2 {
		var sb strings.Builder

		for x := 0; x < int(w); x++ {
			topY := y
			bottomY := y + 1

			var topPixel []byte
			if topY >= 0 {
				topIdx := topY*int(w) + x
				if topIdx < len(pixels) {
					topPixel = pixels[topIdx]
				} else {
					topPixel = []byte{0, 0, 0, 0}
				}
			}

			var bottomPixel []byte
			bottomIdx := bottomY * int(w) + x
			if bottomIdx < len(pixels) {
				bottomPixel = pixels[bottomIdx]
			} else {
				bottomPixel = []byte{0, 0, 0, 0}
			}

			topTransparent := topPixel == nil || topPixel[3] < 32
			bottomTransparent := bottomPixel == nil || bottomPixel[3] < 32

			switch {
			case !topTransparent && !bottomTransparent:
				// Both opaque: BG=top color, FG=bottom color, lower half block
				if useColor {
					sb.WriteString(TrueColorBg(topPixel[0], topPixel[1], topPixel[2]))
					sb.WriteString(TrueColorFg(bottomPixel[0], bottomPixel[1], bottomPixel[2]))
				}
				sb.WriteRune(lowerHalf)

			case !topTransparent && bottomTransparent:
				// Top opaque, bottom transparent: FG=top color, upper half block
				if useColor {
					sb.WriteString(AnsiReset())
					sb.WriteString(TrueColorFg(topPixel[0], topPixel[1], topPixel[2]))
				}
				sb.WriteRune(upperHalf)

			case topTransparent && !bottomTransparent:
				// Top transparent, bottom opaque: FG=bottom color, lower half block
				if useColor {
					sb.WriteString(AnsiReset())
					sb.WriteString(TrueColorFg(bottomPixel[0], bottomPixel[1], bottomPixel[2]))
				}
				sb.WriteRune(lowerHalf)

			default:
				// Both transparent: reset + space
				if useColor {
					sb.WriteString(AnsiReset())
				}
				sb.WriteByte(' ')
			}
		}

		if useColor {
			sb.WriteString(AnsiReset())
			sb.WriteString("\x1b[K")
		}
		sb.WriteByte('\n')

		if _, err := fmt.Fprint(out, sb.String()); err != nil {
			return err
		}
	}

	// Trailing blank line (matching PS Write-Host '' after the loop)
	_, err := fmt.Fprintln(out)
	return err
}

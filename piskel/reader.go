package piskel

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"strings"
)

type piskelFile struct {
	Piskel struct {
		Name   string   `json:"name"`
		Width  int      `json:"width"`
		Height int      `json:"height"`
		Layers []string `json:"layers"`
	} `json:"piskel"`
}

type piskelLayer struct {
	Name       string        `json:"name"`
	FrameCount int           `json:"frameCount"`
	Chunks     []piskelChunk `json:"chunks"`
}

type piskelChunk struct {
	Base64PNG string `json:"base64PNG"`
}

// ParseDimensions extracts canvas width and height from .piskel JSON data.
func ParseDimensions(data []byte) (width, height int, err error) {
	var f piskelFile
	if err := json.Unmarshal(data, &f); err != nil {
		return 0, 0, fmt.Errorf("parse piskel JSON: %w", err)
	}
	return f.Piskel.Width, f.Piskel.Height, nil
}

// ReadLayers decodes each layer's embedded PNG and returns flat RGBA byte slices.
// Layers whose decoded dimensions don't match w x h are silently skipped.
func ReadLayers(data []byte, w, h int) ([][]byte, error) {
	expected := w * h * 4

	var f piskelFile
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse piskel JSON: %w", err)
	}

	var layers [][]byte
	for _, layerJSON := range f.Piskel.Layers {
		var l piskelLayer
		if err := json.Unmarshal([]byte(layerJSON), &l); err != nil {
			continue
		}
		if len(l.Chunks) == 0 {
			continue
		}

		b64 := l.Chunks[0].Base64PNG
		if idx := strings.Index(b64, ","); idx >= 0 {
			b64 = b64[idx+1:]
		}

		pngBytes, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			continue
		}

		img, _, err := image.Decode(bytes.NewReader(pngBytes))
		if err != nil {
			continue
		}

		rgba := imageToRGBA(img)
		if len(rgba) == expected {
			layers = append(layers, rgba)
		}
	}

	return layers, nil
}

// imageToRGBA converts an image.Image to a flat RGBA byte slice.
func imageToRGBA(img image.Image) []byte {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	result := make([]byte, w*h*4)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.At(bounds.Min.X+x, bounds.Min.Y+y)
			r, g, b, a := toNRGBA(c)
			i := (y*w + x) * 4
			result[i] = r
			result[i+1] = g
			result[i+2] = b
			result[i+3] = a
		}
	}
	return result
}

// toNRGBA converts a color to non-premultiplied RGBA bytes.
func toNRGBA(c color.Color) (r, g, b, a uint8) {
	nc := color.NRGBAModel.Convert(c).(color.NRGBA)
	return nc.R, nc.G, nc.B, nc.A
}

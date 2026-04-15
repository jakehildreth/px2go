package piskel_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
)

// buildMinimalPNG creates a valid PNG file from flat RGBA pixel data using Go stdlib.
func buildMinimalPNG(width, height int, rgba []byte) []byte {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := (y*width + x) * 4
			img.SetNRGBA(x, y, color.NRGBA{
				R: rgba[i],
				G: rgba[i+1],
				B: rgba[i+2],
				A: rgba[i+3],
			})
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

// buildMinimalPiskelJSON constructs valid .piskel JSON with embedded base64 PNG layers.
// Note: piskel layers field is an array of JSON strings (not objects).
func buildMinimalPiskelJSON(name string, width, height int, layerNames []string, layerPNGs [][]byte) []byte {
	type chunk struct {
		Layout    [][]int `json:"layout"`
		Base64PNG string  `json:"base64PNG"`
	}
	type layer struct {
		Name       string  `json:"name"`
		Opacity    int     `json:"opacity"`
		FrameCount int     `json:"frameCount"`
		Chunks     []chunk `json:"chunks"`
	}

	layerStrings := make([]string, len(layerNames))
	for i, n := range layerNames {
		b64 := base64.StdEncoding.EncodeToString(layerPNGs[i])
		l := layer{
			Name:       n,
			Opacity:    1,
			FrameCount: 1,
			Chunks: []chunk{{
				Layout:    [][]int{{0}},
				Base64PNG: "data:image/png;base64," + b64,
			}},
		}
		j, _ := json.Marshal(l)
		layerStrings[i] = string(j)
	}

	outer := struct {
		ModelVersion int `json:"modelVersion"`
		Piskel       struct {
			Name        string   `json:"name"`
			Description string   `json:"description"`
			FPS         int      `json:"fps"`
			Height      int      `json:"height"`
			Width       int      `json:"width"`
			Layers      []string `json:"layers"`
		} `json:"piskel"`
	}{
		ModelVersion: 2,
	}
	outer.Piskel.Name = name
	outer.Piskel.Description = ""
	outer.Piskel.FPS = 12
	outer.Piskel.Height = height
	outer.Piskel.Width = width
	outer.Piskel.Layers = layerStrings

	result, _ := json.Marshal(outer)
	return result
}

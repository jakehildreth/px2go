package px_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jakehildreth/px2go/px"
)

func TestParseDimensions(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "Gilmourltd.px"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	w, h, err := px.ParseDimensions(data)
	if err != nil {
		t.Fatalf("ParseDimensions: %v", err)
	}
	if w != 48 {
		t.Errorf("width: got %d, want 48", w)
	}
	if h != 9 {
		t.Errorf("height: got %d, want 9", h)
	}
}

func TestParseDimensions_TooShort(t *testing.T) {
	_, _, err := px.ParseDimensions([]byte{0, 1, 2})
	if err == nil {
		t.Error("expected error for too-short data, got nil")
	}
}

func TestReadLayers(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "Gilmourltd.px"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	layers, err := px.ReadLayers(data, 48, 9)
	if err != nil {
		t.Fatalf("ReadLayers: %v", err)
	}
	if len(layers) != 1 {
		t.Fatalf("layer count: got %d, want 1", len(layers))
	}
	if len(layers[0]) != 48*9*4 {
		t.Errorf("layer size: got %d, want %d", len(layers[0]), 48*9*4)
	}
}

func TestMergeLayers_SingleLayer(t *testing.T) {
	const pixelCount = 4
	layer := make([]byte, pixelCount*4)
	for i := 0; i < pixelCount; i++ {
		layer[i*4] = 255
		layer[i*4+1] = 0
		layer[i*4+2] = 0
		layer[i*4+3] = 255
	}

	pixels := px.MergeLayers([][]byte{layer}, 2, 2)
	if len(pixels) != pixelCount {
		t.Fatalf("pixel count: got %d, want %d", len(pixels), pixelCount)
	}
	for i, p := range pixels {
		if p[0] != 255 || p[1] != 0 || p[2] != 0 || p[3] != 255 {
			t.Errorf("pixel %d: got RGBA(%d,%d,%d,%d), want (255,0,0,255)", i, p[0], p[1], p[2], p[3])
		}
	}
}

func TestMergeLayers_TopLayerWinsWhenOpaque(t *testing.T) {
	bottom := []byte{255, 0, 0, 255} // red opaque
	top := []byte{0, 0, 255, 200}    // blue semi-opaque (alpha > 0)

	pixels := px.MergeLayers([][]byte{top, bottom}, 1, 1)
	if len(pixels) != 1 {
		t.Fatalf("expected 1 pixel, got %d", len(pixels))
	}
	if pixels[0][2] != 255 || pixels[0][0] != 0 {
		t.Errorf("expected blue from top layer; got R=%d B=%d", pixels[0][0], pixels[0][2])
	}
}

func TestMergeLayers_TransparentTopRevealBottom(t *testing.T) {
	bottom := []byte{255, 0, 0, 255} // red opaque
	top := []byte{0, 0, 255, 0}      // blue fully transparent (alpha == 0)

	pixels := px.MergeLayers([][]byte{top, bottom}, 1, 1)
	if pixels[0][0] != 255 || pixels[0][2] != 0 {
		t.Errorf("expected red from bottom layer; got R=%d B=%d", pixels[0][0], pixels[0][2])
	}
}

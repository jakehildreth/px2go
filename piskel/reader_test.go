package piskel_test

import (
	"testing"

	"github.com/jakehildreth/px2go/piskel"
)

func TestParseDimensions(t *testing.T) {
	rgba := make([]byte, 2*2*4)
	pngBytes := buildMinimalPNG(2, 2, rgba)
	data := buildMinimalPiskelJSON("Test", 2, 2, []string{"Layer 1"}, [][]byte{pngBytes})

	w, h, err := piskel.ParseDimensions(data)
	if err != nil {
		t.Fatalf("ParseDimensions: %v", err)
	}
	if w != 2 {
		t.Errorf("width: got %d, want 2", w)
	}
	if h != 2 {
		t.Errorf("height: got %d, want 2", h)
	}
}

func TestReadLayers_SingleLayer_2x2(t *testing.T) {
	rgba := []byte{
		255, 0, 0, 255, // red
		0, 255, 0, 255, // green
		0, 0, 255, 255, // blue
		128, 128, 128, 255, // gray
	}
	pngBytes := buildMinimalPNG(2, 2, rgba)
	data := buildMinimalPiskelJSON("Test", 2, 2, []string{"Layer 1"}, [][]byte{pngBytes})

	layers, err := piskel.ReadLayers(data, 2, 2)
	if err != nil {
		t.Fatalf("ReadLayers: %v", err)
	}
	if len(layers) != 1 {
		t.Fatalf("layer count: got %d, want 1", len(layers))
	}
	if len(layers[0]) != 16 {
		t.Fatalf("layer byte count: got %d, want 16", len(layers[0]))
	}

	// pixel 0: red
	if layers[0][0] != 255 || layers[0][1] != 0 || layers[0][2] != 0 || layers[0][3] != 255 {
		t.Errorf("pixel 0: got RGBA(%d,%d,%d,%d), want (255,0,0,255)",
			layers[0][0], layers[0][1], layers[0][2], layers[0][3])
	}
	// pixel 1: green
	if layers[0][4] != 0 || layers[0][5] != 255 || layers[0][6] != 0 || layers[0][7] != 255 {
		t.Errorf("pixel 1: got RGBA(%d,%d,%d,%d), want (0,255,0,255)",
			layers[0][4], layers[0][5], layers[0][6], layers[0][7])
	}
	// pixel 2: blue
	if layers[0][8] != 0 || layers[0][9] != 0 || layers[0][10] != 255 || layers[0][11] != 255 {
		t.Errorf("pixel 2: got RGBA(%d,%d,%d,%d), want (0,0,255,255)",
			layers[0][8], layers[0][9], layers[0][10], layers[0][11])
	}
	// pixel 3: gray
	if layers[0][12] != 128 || layers[0][13] != 128 || layers[0][14] != 128 || layers[0][15] != 255 {
		t.Errorf("pixel 3: got RGBA(%d,%d,%d,%d), want (128,128,128,255)",
			layers[0][12], layers[0][13], layers[0][14], layers[0][15])
	}
}

func TestReadLayers_MultiLayer(t *testing.T) {
	layer1RGBA := []byte{
		255, 0, 0, 255, // red at (0,0)
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}
	layer2RGBA := []byte{
		0, 0, 0, 0,
		0, 255, 0, 255, // green at (1,0)
		0, 0, 0, 0,
		0, 0, 0, 0,
	}
	png1 := buildMinimalPNG(2, 2, layer1RGBA)
	png2 := buildMinimalPNG(2, 2, layer2RGBA)
	data := buildMinimalPiskelJSON("Test", 2, 2, []string{"Bottom", "Top"}, [][]byte{png1, png2})

	layers, err := piskel.ReadLayers(data, 2, 2)
	if err != nil {
		t.Fatalf("ReadLayers: %v", err)
	}
	if len(layers) != 2 {
		t.Fatalf("layer count: got %d, want 2", len(layers))
	}

	// first layer has red at (0,0)
	if layers[0][0] != 255 || layers[0][1] != 0 {
		t.Errorf("layer 0 pixel 0: got R=%d G=%d, want R=255 G=0", layers[0][0], layers[0][1])
	}
	// second layer has green at (1,0)
	if layers[1][4] != 0 || layers[1][5] != 255 {
		t.Errorf("layer 1 pixel 1: got R=%d G=%d, want R=0 G=255", layers[1][4], layers[1][5])
	}
}

func TestReadLayers_DimensionMismatch(t *testing.T) {
	// create a 1x1 PNG but claim the piskel is 2x2
	smallRGBA := []byte{255, 0, 0, 255}
	smallPNG := buildMinimalPNG(1, 1, smallRGBA)
	data := buildMinimalPiskelJSON("Test", 2, 2, []string{"Mismatched"}, [][]byte{smallPNG})

	layers, err := piskel.ReadLayers(data, 2, 2)
	if err != nil {
		t.Fatalf("ReadLayers: unexpected error %v", err)
	}
	if len(layers) != 0 {
		t.Errorf("layer count: got %d, want 0 (mismatched layer should be skipped)", len(layers))
	}
}

package aseprite_test

import (
	"encoding/binary"
	"testing"

	"github.com/jakehildreth/px2go/aseprite"
)

func TestParseDimensions_4x3(t *testing.T) {
	rgba := make([]byte, 4*3*4) // 4x3 RGBA, all zeros
	data := buildMinimalAseFile(4, 3, []string{"Layer 1"}, [][]byte{rgba})

	w, h, depth, err := aseprite.ParseDimensions(data)
	if err != nil {
		t.Fatalf("ParseDimensions: %v", err)
	}
	if w != 4 {
		t.Errorf("width: got %d, want 4", w)
	}
	if h != 3 {
		t.Errorf("height: got %d, want 3", h)
	}
	if depth != 32 {
		t.Errorf("colorDepth: got %d, want 32", depth)
	}
}

func TestParseDimensions_16x16(t *testing.T) {
	rgba := make([]byte, 16*16*4)
	data := buildMinimalAseFile(16, 16, []string{"BG"}, [][]byte{rgba})

	w, h, _, err := aseprite.ParseDimensions(data)
	if err != nil {
		t.Fatalf("ParseDimensions: %v", err)
	}
	if w != 16 {
		t.Errorf("width: got %d, want 16", w)
	}
	if h != 16 {
		t.Errorf("height: got %d, want 16", h)
	}
}

func TestParseDimensions_BadMagic(t *testing.T) {
	data := make([]byte, 128) // all zeros, magic will be 0x0000
	_, _, _, err := aseprite.ParseDimensions(data)
	if err == nil {
		t.Error("expected error for bad magic number, got nil")
	}
}

func TestParseDimensions_TooShort(t *testing.T) {
	data := []byte{0xE0, 0xA5}
	_, _, _, err := aseprite.ParseDimensions(data)
	if err == nil {
		t.Error("expected error for too-short data, got nil")
	}
}

func TestReadLayers_SingleLayer_2x2(t *testing.T) {
	testRGBA := []byte{
		255, 0, 0, 255, // red
		0, 255, 0, 255, // green
		0, 0, 255, 255, // blue
		128, 128, 128, 255, // gray
	}
	data := buildMinimalAseFile(2, 2, []string{"Layer 1"}, [][]byte{testRGBA})

	layers, err := aseprite.ReadLayers(data, 2, 2)
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
	layer1 := []byte{
		255, 0, 0, 255, // red at (0,0)
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}
	layer2 := []byte{
		0, 0, 0, 0,
		0, 255, 0, 255, // green at (1,0)
		0, 0, 0, 0,
		0, 0, 0, 0,
	}
	data := buildMinimalAseFile(2, 2, []string{"Bottom", "Top"}, [][]byte{layer1, layer2})

	layers, err := aseprite.ReadLayers(data, 2, 2)
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

func TestReadLayers_NoCelData(t *testing.T) {
	// header-only file with 0 frames
	data := make([]byte, 128)
	binary.LittleEndian.PutUint32(data[0:4], 128)
	binary.LittleEndian.PutUint16(data[4:6], 0xA5E0)
	binary.LittleEndian.PutUint16(data[6:8], 0)    // 0 frames
	binary.LittleEndian.PutUint16(data[8:10], 2)   // width
	binary.LittleEndian.PutUint16(data[10:12], 2)  // height
	binary.LittleEndian.PutUint16(data[12:14], 32) // 32bpp

	layers, err := aseprite.ReadLayers(data, 2, 2)
	if err != nil {
		t.Fatalf("ReadLayers: unexpected error %v", err)
	}
	if len(layers) != 0 {
		t.Errorf("layer count: got %d, want 0", len(layers))
	}
}

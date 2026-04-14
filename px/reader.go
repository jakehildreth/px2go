package px

import (
	"encoding/binary"
	"fmt"
)

// ParseDimensions extracts canvas width and height from a Pixquare .px file.
// Format: byte 8 = ID string length N; UInt32LE width at 64+N; UInt32LE height at 64+N+4.
func ParseDimensions(data []byte) (width, height uint32, err error) {
	if len(data) < 9 {
		return 0, 0, fmt.Errorf("data too short (%d bytes)", len(data))
	}
	idLen := int(data[8])
	sizeOffset := 64 + idLen
	if len(data) < sizeOffset+8 {
		return 0, 0, fmt.Errorf("data too short for dimensions at offset %d", sizeOffset)
	}
	width = binary.LittleEndian.Uint32(data[sizeOffset:])
	height = binary.LittleEndian.Uint32(data[sizeOffset+4:])
	return width, height, nil
}

// ReadLayers finds and decompresses all valid zlib-encoded pixel layers in data.
// Each valid layer must decompress to exactly w*h*4 bytes (RGBA).
func ReadLayers(data []byte, w, h uint32) ([][]byte, error) {
	expected := int(w) * int(h) * 4
	offsets := FindZlibHeaders(data)
	if len(offsets) == 0 {
		return nil, fmt.Errorf("no zlib headers found")
	}

	var layers [][]byte
	for _, off := range offsets {
		expanded, err := ExpandZlibData(data, off)
		if err != nil {
			continue
		}
		if len(expanded) == expected {
			layers = append(layers, expanded)
		}
	}
	if len(layers) == 0 {
		return nil, fmt.Errorf("no valid layers found (expected %d bytes each)", expected)
	}
	return layers, nil
}

// MergeLayers composites layers bottom-to-top into a flat pixel slice.
// layers[0] is the topmost; layers[len-1] is the bottom.
// A pixel from a higher layer replaces the pixel below if its alpha > 0.
func MergeLayers(layers [][]byte, w, h uint32) [][]byte {
	pixelCount := int(w) * int(h)
	result := make([][]byte, pixelCount)

	bottom := layers[len(layers)-1]
	for i := 0; i < pixelCount; i++ {
		pixel := make([]byte, 4)
		copy(pixel, bottom[i*4:i*4+4])
		result[i] = pixel
	}

	for l := len(layers) - 2; l >= 0; l-- {
		layer := layers[l]
		for i := 0; i < pixelCount; i++ {
			if layer[i*4+3] > 0 {
				copy(result[i], layer[i*4:i*4+4])
			}
		}
	}

	return result
}

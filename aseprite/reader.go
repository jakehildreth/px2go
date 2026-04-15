package aseprite

import (
	"encoding/binary"
	"fmt"

	"github.com/jakehildreth/px2go/px"
)

const (
	headerSize  = 128
	magicNumber = 0xA5E0
	frameMagic  = 0xF1FA
	chunkLayer  = 0x2004
	chunkCel    = 0x2005
	celTypeComp = 2
)

// ParseDimensions extracts canvas width, height, and color depth from an
// Aseprite .ase file header.
func ParseDimensions(data []byte) (width, height, colorDepth uint16, err error) {
	if len(data) < headerSize {
		return 0, 0, 0, fmt.Errorf("data too short for Aseprite header (got %d bytes, need %d)", len(data), headerSize)
	}
	magic := binary.LittleEndian.Uint16(data[4:6])
	if magic != magicNumber {
		return 0, 0, 0, fmt.Errorf("invalid Aseprite magic number: 0x%04X (expected 0x%04X)", magic, magicNumber)
	}
	width = binary.LittleEndian.Uint16(data[8:10])
	height = binary.LittleEndian.Uint16(data[10:12])
	colorDepth = binary.LittleEndian.Uint16(data[12:14])
	return width, height, colorDepth, nil
}

// ReadLayers parses the first frame of an Aseprite file and returns
// decompressed RGBA data for each cel. Only processes compressed image
// cels (type 2) that match the expected dimensions.
func ReadLayers(data []byte, w, h uint16) ([][]byte, error) {
	expected := int(w) * int(h) * 4

	if len(data) < headerSize {
		return nil, nil
	}

	frameCount := binary.LittleEndian.Uint16(data[6:8])
	if frameCount == 0 {
		return nil, nil
	}

	// only process first frame
	offset := headerSize
	if offset+16 > len(data) {
		return nil, nil
	}

	fMagic := binary.LittleEndian.Uint16(data[offset+4 : offset+6])
	if fMagic != frameMagic {
		return nil, fmt.Errorf("invalid frame magic: 0x%04X", fMagic)
	}

	oldChunkCount := binary.LittleEndian.Uint16(data[offset+6 : offset+8])
	newChunkCount := binary.LittleEndian.Uint32(data[offset+12 : offset+16])
	chunkCount := int(oldChunkCount)
	if newChunkCount != 0 {
		chunkCount = int(newChunkCount)
	}

	chunkOffset := offset + 16
	var layers [][]byte

	for ci := 0; ci < chunkCount; ci++ {
		if chunkOffset+6 > len(data) {
			break
		}

		chunkSize := int(binary.LittleEndian.Uint32(data[chunkOffset : chunkOffset+4]))
		chunkType := binary.LittleEndian.Uint16(data[chunkOffset+4 : chunkOffset+6])

		if chunkType == chunkCel {
			celData := chunkOffset + 6
			if celData+20 > len(data) {
				chunkOffset += chunkSize
				continue
			}

			celType := binary.LittleEndian.Uint16(data[celData+7 : celData+9])
			if celType == celTypeComp {
				celW := binary.LittleEndian.Uint16(data[celData+16 : celData+18])
				celH := binary.LittleEndian.Uint16(data[celData+18 : celData+20])

				if celW == w && celH == h {
					compStart := celData + 20
					expanded, err := px.ExpandZlibData(data, compStart)
					if err == nil && len(expanded) == expected {
						layers = append(layers, expanded)
					}
				}
			}
		}

		chunkOffset += chunkSize
	}

	return layers, nil
}

package aseprite_test

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
)

// buildMinimalAseFile constructs a valid .ase binary from known pixel data.
// It creates a 128-byte header, one frame, and one layer chunk (0x2004) +
// one compressed cel chunk (0x2005 type 2) per layer.
func buildMinimalAseFile(width, height uint16, layerNames []string, layerRGBA [][]byte) []byte {
	var frame bytes.Buffer
	chunkCount := uint16(len(layerNames) * 2)

	// frame header placeholder (16 bytes)
	binary.Write(&frame, binary.LittleEndian, uint32(0))      // frame size placeholder
	binary.Write(&frame, binary.LittleEndian, uint16(0xF1FA)) // magic
	binary.Write(&frame, binary.LittleEndian, chunkCount)     // old chunk count
	binary.Write(&frame, binary.LittleEndian, uint16(100))    // frame duration
	frame.Write([]byte{0, 0})                                 // future
	binary.Write(&frame, binary.LittleEndian, uint32(0))      // new chunk count (0 = use old)

	for li := 0; li < len(layerNames); li++ {
		nameBytes := []byte(layerNames[li])

		// layer chunk (0x2004)
		var layerChunk bytes.Buffer
		binary.Write(&layerChunk, binary.LittleEndian, uint32(0))      // size placeholder
		binary.Write(&layerChunk, binary.LittleEndian, uint16(0x2004)) // chunk type
		binary.Write(&layerChunk, binary.LittleEndian, uint16(1))      // flags: visible
		binary.Write(&layerChunk, binary.LittleEndian, uint16(0))      // layer type: normal
		binary.Write(&layerChunk, binary.LittleEndian, uint16(0))      // child level
		binary.Write(&layerChunk, binary.LittleEndian, uint16(0))      // default width
		binary.Write(&layerChunk, binary.LittleEndian, uint16(0))      // default height
		binary.Write(&layerChunk, binary.LittleEndian, uint16(0))      // blend mode
		layerChunk.WriteByte(255)                                      // opacity
		layerChunk.Write([]byte{0, 0, 0})                              // future
		binary.Write(&layerChunk, binary.LittleEndian, uint16(len(nameBytes)))
		layerChunk.Write(nameBytes)

		// patch chunk size
		chunkBytes := layerChunk.Bytes()
		binary.LittleEndian.PutUint32(chunkBytes[0:4], uint32(len(chunkBytes)))
		frame.Write(chunkBytes)

		// cel chunk (0x2005) - compressed image
		var celChunk bytes.Buffer
		binary.Write(&celChunk, binary.LittleEndian, uint32(0))      // size placeholder
		binary.Write(&celChunk, binary.LittleEndian, uint16(0x2005)) // chunk type
		binary.Write(&celChunk, binary.LittleEndian, uint16(li))     // layer index
		binary.Write(&celChunk, binary.LittleEndian, int16(0))       // x position
		binary.Write(&celChunk, binary.LittleEndian, int16(0))       // y position
		celChunk.WriteByte(255)                                      // opacity
		binary.Write(&celChunk, binary.LittleEndian, uint16(2))      // cel type: compressed
		binary.Write(&celChunk, binary.LittleEndian, int16(0))       // z-index
		celChunk.Write([]byte{0, 0, 0, 0, 0})                        // future
		binary.Write(&celChunk, binary.LittleEndian, width)          // cel width
		binary.Write(&celChunk, binary.LittleEndian, height)         // cel height

		// zlib header + deflate compressed RGBA data
		celChunk.Write([]byte{0x78, 0x9C})
		var deflated bytes.Buffer
		w, _ := flate.NewWriter(&deflated, flate.DefaultCompression)
		w.Write(layerRGBA[li])
		w.Close()
		celChunk.Write(deflated.Bytes())

		// patch chunk size
		celBytes := celChunk.Bytes()
		binary.LittleEndian.PutUint32(celBytes[0:4], uint32(len(celBytes)))
		frame.Write(celBytes)
	}

	// patch frame size
	frameBytes := frame.Bytes()
	binary.LittleEndian.PutUint32(frameBytes[0:4], uint32(len(frameBytes)))

	// 128-byte header
	header := make([]byte, 128)
	fileSize := uint32(128 + len(frameBytes))
	binary.LittleEndian.PutUint32(header[0:4], fileSize)
	binary.LittleEndian.PutUint16(header[4:6], 0xA5E0) // magic
	binary.LittleEndian.PutUint16(header[6:8], 1)      // 1 frame
	binary.LittleEndian.PutUint16(header[8:10], width)
	binary.LittleEndian.PutUint16(header[10:12], height)
	binary.LittleEndian.PutUint16(header[12:14], 32)  // 32bpp RGBA
	binary.LittleEndian.PutUint32(header[14:18], 1)   // flags: layer opacity valid
	binary.LittleEndian.PutUint16(header[18:20], 100) // speed
	header[34] = 1                                    // pixel width
	header[35] = 1                                    // pixel height
	binary.LittleEndian.PutUint16(header[40:42], 16)  // grid width
	binary.LittleEndian.PutUint16(header[42:44], 16)  // grid height

	result := make([]byte, 0, len(header)+len(frameBytes))
	result = append(result, header...)
	result = append(result, frameBytes...)
	return result
}

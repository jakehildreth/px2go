package px

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
)

// zlibMagic holds all valid zlib 2-byte header prefixes.
var zlibMagic = [][2]byte{
	{0x78, 0x9C},
	{0x78, 0xDA},
	{0x78, 0x01},
	{0x78, 0x5E},
}

// FindZlibHeaders scans data and returns the byte offset of every valid zlib header.
func FindZlibHeaders(data []byte) []int {
	var offsets []int
	for i := 0; i < len(data)-1; i++ {
		for _, magic := range zlibMagic {
			if data[i] == magic[0] && data[i+1] == magic[1] {
				offsets = append(offsets, i)
				break
			}
		}
	}
	return offsets
}

// ExpandZlibData decompresses a zlib stream at offset in data.
// It skips the 2-byte zlib header and decompresses via raw deflate.
func ExpandZlibData(data []byte, offset int) ([]byte, error) {
	if offset+2 > len(data) {
		return nil, fmt.Errorf("offset %d out of range (data len %d)", offset, len(data))
	}
	r := flate.NewReader(bytes.NewReader(data[offset+2:]))
	defer r.Close()
	return io.ReadAll(r)
}

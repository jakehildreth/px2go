package px_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jakehildreth/px2go/px"
)

func TestFindZlibHeaders(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "Gilmourltd.px"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	offsets := px.FindZlibHeaders(data)
	if len(offsets) != 1 {
		t.Fatalf("expected 1 zlib header, got %d: %v", len(offsets), offsets)
	}
	if offsets[0] != 502 {
		t.Errorf("expected offset 502, got %d", offsets[0])
	}
}

func TestExpandZlibData(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "Gilmourltd.px"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	const wantLen = 48 * 9 * 4 // 1728 bytes

	expanded, err := px.ExpandZlibData(data, 502)
	if err != nil {
		t.Fatalf("ExpandZlibData: %v", err)
	}
	if len(expanded) != wantLen {
		t.Errorf("decompressed %d bytes, want %d", len(expanded), wantLen)
	}
}

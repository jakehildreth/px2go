# Plan: Go CLI Port of PX2PS → `px2go`

**TL;DR:** TDD throughout. write failing test → implement to pass → refactor. `io.Writer` param on Render makes it trivially testable without stdout hijacking.

## Decisions
- Separate repo: github.com/jakehildreth/px2go (domain: px2go.jakehildreth.com)
- Module path: `github.com/jakehildreth/px2go`
- Binary name: `px2go`
- Output modes: Display only (ANSI true color to stdout)
- Color: truecolor only; auto-detect via COLORTERM/NO_COLOR/Windows VT API
- No ConsoleColor fallback, no ScriptBlock/Script modes
- No cobra; stdlib `flag` package
- Dependencies: stdlib + `golang.org/x/sys` (Windows VT)
- Render func accepts io.Writer (not hardcoded os.Stdout) for testability
- TDD: write failing test first, implement to pass, refactor

## Project Structure
```
px2go/
├── go.mod                   module: github.com/jakehildreth/px2go, Go 1.22+
├── go.sum
├── main.go
├── px/
│   ├── testdata/
│   │   └── Gilmourltd.px    copied from px2ps/Examples/ as fixture
│   ├── zlib.go
│   ├── zlib_test.go
│   ├── reader.go
│   └── reader_test.go
└── render/
    ├── color.go
    ├── color_test.go
    ├── terminal.go          Render(w, h, pixels, colorMode, io.Writer) error
    └── terminal_test.go
```

## Px File Format
- Bytes 0-63: header (fixed)
- Byte 8: ID string length (N)
- Bytes 64+N: UInt32LE width
- Bytes 64+N+4: UInt32LE height
- Remaining: zlib layers (valid headers: 0x789C, 0x78DA, 0x7801, 0x785E)
- Each decompressed layer: Width × Height × 4 bytes (RGBA)
- Compositing: bottom-to-top, replace if alpha > 0

## Rendering Algorithm
- 2 pixel rows → 1 terminal row
- Odd height: start at Y=-1 (top row = upper half block only)
- Transparency threshold: alpha < 32
- Both opaque: BG=top, FG=bottom, char=▄
- Top opaque, bottom transparent: FG=top, char=▀
- Top transparent, bottom opaque: FG=bottom, char=▄
- Both transparent: reset + space
- End of each row: reset + ESC[K (erase to EOL)

## CLI Interface
```
px2go [flags] <path> [<path>...]
  -color string    auto|truecolor|none (default: auto)
  -v               verbose
```
- directory arg → glob all *.px recursively

## Color Detection (auto)
1. NO_COLOR env set → none
2. Windows: SetConsoleMode VT; fail → none
3. COLORTERM=truecolor|24bit → truecolor
4. Default → truecolor

## Implementation Steps (TDD: RED → GREEN → REFACTOR each phase)

### Phase 1: zlib
1. `px/zlib_test.go`: TestFindZlibHeaders (known offsets in Gilmourltd.px fixture), TestExpandZlibData (decompressed len = W×H×4)
2. `px/zlib.go`: implement FindZlibHeaders + ExpandZlibData

### Phase 2: px parsing
3. `px/reader_test.go`: TestParseDimensions (expected W/H from fixture), TestReadLayers (count + sizes), TestMergeLayers (single-layer passthrough; multi-layer alpha composite with synthetic data)
4. `px/reader.go`: implement ParseDimensions, ReadLayers, MergeLayers

### Phase 3: color helpers
5. `render/color_test.go`: TestTrueColorFg/Bg (exact escape strings), TestAnsiReset, TestDetectColorMode (NO_COLOR override)
6. `render/color.go`: implement

### Phase 4: rendering
7. `render/terminal_test.go`: TestRender with io.Writer capture; cover even/odd height, all 4 transparency cases, reset+erase per row
8. `render/terminal.go`: implement Render(w, h uint32, pixels [][]byte, colorMode string, out io.Writer) error

### Phase 5: CLI wiring
9. `main.go`: file/dir resolution, loop → parse → render; wire out=os.Stdout

### Phase 6: platform
10. Windows VT enablement (runtime.GOOS check, golang.org/x/sys/windows.SetConsoleMode)
11. go.mod finalized with all deps

## Verification
1. `go test ./...` all green
2. `go build ./...` + `go vet ./...` clean
3. `./px2go px/testdata/Gilmourltd.px` renders correctly in terminal
4. `./px2go px/testdata/` processes all `.px` files

# px2go

A cross-platform CLI that converts [Pixquare](https://pixquare.app) .px files to terminal pixel graphics using ANSI True Color.

## Installation

From source:
```bash
git clone https://github.com/jakehildreth/px2go.git
cd px2go
go build -o px2go .
```

Pre-built binaries for macOS, Linux, and Windows will be available on the [releases page](https://github.com/jakehildreth/px2go/releases).

## Quick Start

```bash
# Render a .px file
px2go image.px

# Render all .px files in a directory
px2go ./artworks/

# Render multiple files
px2go logo.px banner.px

# Verbose output (shows dimensions and layer count)
px2go -v image.px

# Force a specific color mode
px2go -color truecolor image.px
px2go -color none image.px
```

## Features

- [x] Renders .px files directly in the terminal using ANSI True Color
- [x] Supports single-layer and multi-layer .px files
- [x] Automatic layer compositing and transparency handling
- [x] Cross-platform: macOS, Linux, Windows
- [x] Works in bash, zsh, fish, PowerShell, and any modern shell
- [x] Zero runtime dependencies

## Requirements

- A terminal with ANSI True Color support (iTerm2, Windows Terminal, kitty, most modern terminals)
- Go 1.22+ to build from source

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-color` | `auto` | Color mode: `auto`, `truecolor`, or `none` |
| `-v` | `false` | Verbose output (dimensions, layer count) |

**Color mode auto-detection:** respects `NO_COLOR` env var; checks `COLORTERM=truecolor\|24bit`; enables Windows VT processing automatically.

## How It Works

px2go reads Pixquare .px files, decompresses the zlib-encoded layer data, composites multiple layers if present, and renders the final image using Unicode half-block characters (▄ ▀) with ANSI True Color escape sequences. Each terminal row represents two pixel rows.

## License

MIT License w/Commons Clause - see [LICENSE](LICENSE) file for details.

---

Made with 💜 by [Jake Hildreth](https://jakehildreth.com)


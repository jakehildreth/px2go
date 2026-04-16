# px2go

A cross-platform CLI that renders pixel art in the terminal using ANSI True Color. Supports [Pixquare](https://pixquare.app) (.px), [Aseprite](https://www.aseprite.org) (.ase/.aseprite), and [Piskel](https://www.piskelapp.com) (.piskel) formats.

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
# Render a pixel art file (any supported format)
px2go image.px
px2go sprite.aseprite
px2go animation.piskel

# Render all supported files in a directory
px2go ./artworks/

# Render multiple files
px2go logo.px banner.ase sprite.piskel

# Verbose output (shows dimensions and layer count)
px2go -v image.px

# Force a specific color mode
px2go -color truecolor image.px
px2go -color none image.px
```

## Features

- [x] Renders .px, .ase/.aseprite, and .piskel files directly in the terminal using ANSI True Color
- [x] Supports single-layer and multi-layer images across all formats
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

px2go detects the file format by extension and parses accordingly:

| Format | Extensions | Parser |
|--------|------------|--------|
| Pixquare | `.px` | Zlib header scan + deflate decompression |
| Aseprite | `.ase`, `.aseprite` | Binary header/frame/chunk parsing (first frame only) |
| Piskel | `.piskel` | JSON + embedded base64 PNG decoding (first frame only) |

After parsing, all formats share the same pipeline: layer compositing, transparency handling, and rendering via Unicode half-block characters (▄ ▀) with ANSI True Color escape sequences. Each terminal row represents two pixel rows.

## License

MIT License w/Commons Clause - see [LICENSE](LICENSE) file for details.

---

Made with 💜 by [Jake Hildreth](https://jakehildreth.com)


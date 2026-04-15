package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jakehildreth/px2go/aseprite"
	"github.com/jakehildreth/px2go/piskel"
	"github.com/jakehildreth/px2go/px"
	"github.com/jakehildreth/px2go/render"
)

func main() {
	color := flag.String("color", "auto", "color mode: auto|truecolor|none")
	verbose := flag.Bool("v", false, "verbose output")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: px2go [flags] <path> [<path>...]")
		fmt.Fprintln(os.Stderr, "flags:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	colorMode := *color
	if colorMode == "auto" {
		colorMode = render.DetectColorMode()
	}

	var paths []string
	for _, arg := range flag.Args() {
		info, err := os.Stat(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "px2go: %v\n", err)
			os.Exit(1)
		}
		if info.IsDir() {
			err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() {
					ext := strings.ToLower(filepath.Ext(path))
					switch ext {
					case ".px", ".piskel", ".ase", ".aseprite":
						paths = append(paths, path)
					}
				}
				return nil
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "px2go: %v\n", err)
				os.Exit(1)
			}
		} else {
			paths = append(paths, arg)
		}
	}

	exitCode := 0
	for _, path := range paths {
		if err := processFile(path, colorMode, *verbose); err != nil {
			fmt.Fprintf(os.Stderr, "px2go: %s: %v\n", path, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func processFile(path, colorMode string, verbose bool) error {
	if verbose {
		fmt.Fprintf(os.Stderr, "[i] rendering %s\n", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(path))

	var (
		w, h uint32
		layers [][]byte
	)

	switch ext {
	case ".ase", ".aseprite":
		aw, ah, _, err := aseprite.ParseDimensions(data)
		if err != nil {
			return fmt.Errorf("parse dimensions: %w", err)
		}
		w, h = uint32(aw), uint32(ah)
		layers, err = aseprite.ReadLayers(data, aw, ah)
		if err != nil {
			return fmt.Errorf("read layers: %w", err)
		}

	case ".piskel":
		pw, ph, err := piskel.ParseDimensions(data)
		if err != nil {
			return fmt.Errorf("parse dimensions: %w", err)
		}
		w, h = uint32(pw), uint32(ph)
		layers, err = piskel.ReadLayers(data, pw, ph)
		if err != nil {
			return fmt.Errorf("read layers: %w", err)
		}

	default: // .px
		pxW, pxH, err := px.ParseDimensions(data)
		if err != nil {
			return fmt.Errorf("parse dimensions: %w", err)
		}
		w, h = pxW, pxH
		layers, err = px.ReadLayers(data, pxW, pxH)
		if err != nil {
			return fmt.Errorf("read layers: %w", err)
		}
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "[i] dimensions: %dx%d\n", w, h)
		fmt.Fprintf(os.Stderr, "[i] layers: %d\n", len(layers))
	}

	pixels := px.MergeLayers(layers, w, h)
	return render.Render(w, h, pixels, colorMode, os.Stdout)
}

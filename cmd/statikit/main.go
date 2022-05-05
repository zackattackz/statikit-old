package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zackattackz/azure_static_site_kit/pkg/statikitRenderer"
)

const (
	version     = "v0.1.0"
	usageString = ""

	// Modes of operation
	render  = "render"
	preview = "preview"
	publish = "publish"

	// flag names
	inFlag  = "d"
	outFlag = "o"

	// default flag values
	defaultIn = "."

	// flag descriptions
	descIn  = "unrendered input directory"
	descOut = "rendered output directory"
)

func main() {

	fmt.Printf("statikit %v...\n", version)

	mode := os.Args[1]
	// Set os.Args to the remaining args, so the flag package will ignore the first when parsing flags
	os.Args = os.Args[1:]

	var inDir string
	flag.StringVar(&inDir, inFlag, defaultIn, descIn)

	switch mode {

	case render, publish:
		var outDir string
		defaultOut := os.TempDir() + filepath.FromSlash("/statikitRendered")
		flag.StringVar(&outDir, outFlag, defaultOut, descOut)
		flag.Parse()
		rendererArgs := statikitRenderer.Args{InDir: inDir, OutDir: outDir}
		for _, err := range statikitRenderer.Render(rendererArgs) {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		}

	case preview:

	default:
		fmt.Fprintf(os.Stderr, "%v\n", usageString)
		os.Exit(-1)
	}

}

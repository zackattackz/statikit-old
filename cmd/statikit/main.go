package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zackattackz/azure_static_site_kit/pkg/statikit"
)

const (
	version     = "v0.1.0"
	usageString = "usage: statikit [render | preview | publish] opts"

	// Modes of operation
	render  = "render"
	preview = "preview"
	publish = "publish"

	// flag names
	inFlag    = "d"
	outFlag   = "o"
	forceFlag = "f"

	// default flag values
	defaultIn    = "."
	defaultForce = false

	// flag descriptions
	descIn    = "unrendered input directory"
	descOut   = "rendered output directory"
	descForce = "force output directory removal"
)

func logErrAndExit(err error, code int) {
	fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	os.Exit(code)
}

func warnErase(outDir string) error {
	fmt.Printf("WARNING!! OK to delete everything inside %v y/[n]: ", outDir)
	bio := bufio.NewReader(os.Stdin)
	line, err := bio.ReadString('\n')
	if err != nil {
		return err
	}
	if line == "y\n" || line == "y\r\n" {
		return nil
	}
	return fmt.Errorf("removal of %v not confirmed", outDir)
}

func main() {
	mode := os.Args[1]
	// Set os.Args to the remaining args, so the flag package will ignore the first when parsing flags
	os.Args = os.Args[1:]

	var inDir string
	var outDir string
	var force bool
	defaultOut := filepath.Join(os.TempDir(), "statikitRendered")
	flag.StringVar(&inDir, inFlag, defaultIn, descIn)
	flag.StringVar(&outDir, outFlag, defaultOut, descOut)
	flag.BoolVar(&force, forceFlag, defaultForce, descForce)
	flag.Parse()

	// If invalid mode print usage string and exit
	switch mode {
	case render, preview, publish:
		break
	default:
		logErrAndExit(fmt.Errorf(usageString), 1)
	}

	// If no force flag, ensure user wants to erase.
	if !force {
		if err := warnErase(outDir); err != nil {
			logErrAndExit(err, 1)
		}
	}

	// If we make it here, erase outdir
	if err := os.RemoveAll(outDir); err != nil {
		logErrAndExit(err, 1)
	}

	// Go into mode specific handlers
	switch mode {

	case render:

		rendererArgs := statikit.RendererArgs{InDir: inDir, OutDir: outDir}
		if err := statikit.Render(rendererArgs); err != nil {
			logErrAndExit(err, 1)
		}

	case preview:

	case publish:

	}

}

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zackattackz/azure_static_site_kit/pkg/statikit/previewer"
	"github.com/zackattackz/azure_static_site_kit/pkg/statikit/publisher"
)

const (
	version     = "v0.1.0"
	usageString = "usage: statikit [init | render | preview | publish] opts"

	// Modes of operation
	modeRender  = "render"
	modePreview = "preview"
	modePublish = "publish"
	modeInit    = "init"
)

func logErrAndExit(err error, code int) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(code)
}

func printUsageAndExit() {
	fmt.Fprintln(os.Stderr, usageString)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {

	if len(os.Args) < 2 {
		printUsageAndExit()
	}

	mode := os.Args[1]

	// Set os.Args to the remaining args, so the flag package will ignore the first when parsing flags
	os.Args = os.Args[1:]

	switch mode {

	case modeInit:
		if len(os.Args) < 2 {
			printUsageAndExit()
		}
		outDir := os.Args[1]
		err := initialize(outDir)
		if err != nil {
			// Delete outDir if it was made
			os.RemoveAll(outDir)
			logErrAndExit(err, 1)
		}

	case modeRender:
		const (
			// flags
			inFlag            = "d"
			outFlag           = "o"
			forceFlag         = "f"
			rendererCountFlag = "renderer-count"

			// default flag values
			defaultIn            = "."
			defaultForce         = false
			defaultRendererCount = 20

			// flag descriptions
			descIn            = "unrendered input directory"
			descOut           = "rendered output directory"
			descForce         = "force output directory removal"
			descRendererCount = "how many renderer goroutines to be made"
		)

		a := renderArgs{}

		// Parse all flags into a
		defaultOut := filepath.Join(os.TempDir(), "statikitRendered")
		flag.StringVar(&a.inDir, inFlag, defaultIn, descIn)
		flag.StringVar(&a.outDir, outFlag, defaultOut, descOut)
		flag.BoolVar(&a.force, forceFlag, defaultForce, descForce)
		flag.UintVar(&a.rendererCount, rendererCountFlag, defaultRendererCount, descRendererCount)
		flag.Parse()

		// Clean the in/out dirs
		a.inDir = filepath.Clean(a.inDir)
		a.outDir = filepath.Clean(a.outDir)

		cfg, cfgPath, err := parseConfig(a.inDir)
		if err != nil {
			logErrAndExit(err, 1)
		}
		a.data = cfg.Data

		// Determine cfg file name
		_, a.cfgFileName = filepath.Split(cfgPath)

		err = render(a)
		if err != nil {
			logErrAndExit(err, 1)
		}

	case modePreview:
		if len(os.Args) < 2 {
			printUsageAndExit()
		}
		outDir := os.Args[1]
		err := previewer.Preview(outDir)
		if err != nil {
			logErrAndExit(err, 1)
		}

	case modePublish:
		if len(os.Args) < 2 {
			printUsageAndExit()
		}
		outDir := os.Args[1]
		publisherArgs := publisher.PublisherArgs{Path: outDir}
		publisher.Publish(publisherArgs)
	}
}

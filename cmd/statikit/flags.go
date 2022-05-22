package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/preview"
)

type flagValues struct {
	render  renderArgs
	preview preview.Args
	publish publishArgs
	init    initArgs
}

type modeToFlags map[mode]*flag.FlagSet

func initModeToFlags(mToFs *modeToFlags, fValues *flagValues, cmdName string) {
	// Create all the flag sets
	*mToFs = map[mode]*flag.FlagSet{
		modeInvalid: flag.NewFlagSet(string(modeInvalid), flag.ExitOnError),
		modeRender:  flag.NewFlagSet(string(modeRender), flag.ExitOnError),
		modePreview: flag.NewFlagSet(string(modePreview), flag.ExitOnError),
		modePublish: flag.NewFlagSet(string(modePublish), flag.ExitOnError),
		modeInit:    flag.NewFlagSet(string(modeInit), flag.ExitOnError),
		modeHelp:    flag.NewFlagSet(string(modeHelp), flag.ExitOnError),
	}

	// Init render flags
	func() {
		fs := (*mToFs)[modeRender]
		const (
			// flags
			outFlag           = "o"
			forceFlag         = "f"
			rendererCountFlag = "renderer-count"

			// default flag values
			defaultForce         = false
			defaultRendererCount = 20

			// flag descriptions
			descOut           = "rendered output directory"
			descForce         = "force output directory removal"
			descRendererCount = "how many renderer goroutines to be made"
		)
		defaultOut := filepath.Join(os.TempDir(), "statikitRendered")
		fs.StringVar(&fValues.render.outDir, outFlag, defaultOut, descOut)
		fs.BoolVar(&fValues.render.force, forceFlag, defaultForce, descForce)
		fs.UintVar(&fValues.render.rendererCount, rendererCountFlag, defaultRendererCount, descRendererCount)
		fs.Usage = func() {
			printUsageAndExit(cmdName, modeRender)
		}
	}()

	// func() {
	// 	// Init preview flags
	// }()

}

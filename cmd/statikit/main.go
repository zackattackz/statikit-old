package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/preview"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/configParser"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/initializer"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/previewer"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/schemaParser"
)

func logErrAndExit(err error, code int) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(code)
}

var fValues flagValues

var mToFs modeToFlags

func main() {
	_, cmdName := filepath.Split(os.Args[0])

	initModeToFlags(&mToFs, &fValues, cmdName)

	if len(os.Args) < 2 {
		printUsageAndExit(cmdName, modeInvalid)
	}

	m := mode(os.Args[1])

	remainingArgs := os.Args[2:]

	switch m {

	case modeInit:
		if len(remainingArgs) > 1 {
			printUsageAndExit(cmdName, modeInit)
		}
		outDir := filepath.Clean(remainingArgs[0])
		a := initArgs{path: outDir}
		err := initialize(a)
		if err != nil {
			logErrAndExit(err, 1)
		}

	case modeRender:

		// Parse flags
		flagSet := mToFs[modeRender]
		flagSet.Parse(remainingArgs)

		if flagSet.NArg() > 1 {
			printUsageAndExit(cmdName, modeRender)
		}

		// Initialize renderer args from parsed values
		a := renderArgs{
			outDir:        fValues.render.outDir,
			force:         fValues.render.force,
			rendererCount: fValues.render.rendererCount,
		}

		// Initialize inDir to (optionally) first non-flag arg
		a.inDir = flagSet.Arg(0)
		if a.inDir == "" {
			a.inDir = "."
		}

		// Clean the in/out dirs
		a.inDir = filepath.Clean(a.inDir)
		a.outDir = filepath.Clean(a.outDir)

		// cfg, err := parseConfig(a.inDir)
		// if err != nil {
		// 	logErrAndExit(err, 1)
		// }

		schemaMap := make(schemaParser.Map)
		parser := schemaParser.New(a.inDir)
		err := parser.Parse(&schemaMap)
		if err != nil {
			logErrAndExit(err, 1)
		}
		a.schemaMap = schemaMap

		cfgParser, err := configParser.New(a.inDir)
		if err != nil {
			logErrAndExit(err, 1)
		}
		cfg := configParser.Config{}
		cfgParser.Parse(&cfg)
		a.ignore = cfg.Ignore
		a.ignore = append(a.ignore, initializer.StatikitDirName)

		err = render(a)
		if err != nil {
			logErrAndExit(err, 1)
		}

	case modePreview:
		if len(remainingArgs) > 1 {
			printUsageAndExit(cmdName, modePreview)
		}
		inDir := filepath.Clean(remainingArgs[0])
		p := previewer.New(inDir, ":8080")
		err := preview.Run(p)
		if err != nil {
			logErrAndExit(err, 1)
		}

	case modePublish:
		if len(remainingArgs) > 1 {
			printUsageAndExit(cmdName, modePublish)
		}
		inDir := filepath.Clean(remainingArgs[0])
		a := publishArgs{path: inDir}
		err := publish(a)
		if err != nil {
			logErrAndExit(err, 1)
		}

	case modeHelp:
		if len(remainingArgs) > 1 ||
			len(remainingArgs) < 1 {
			printUsageAndExit(cmdName, modeHelp)
		}
		m := remainingArgs[0]
		if isMode(m) {
			printUsageAndExit(cmdName, mode(m))

		}
		printUsageAndExit(cmdName, modeHelp)

	default:
		printUsageAndExit(cmdName, modeInvalid)
	}
}

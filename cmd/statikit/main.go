package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func logErrAndExit(err error, code int) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(code)
}

var fValues flagValues

var mToFs modeToFlags

func main() {
	_, cmdName := filepath.Split(os.Args[0])

	if len(os.Args) < 2 {
		printUsageAndExit(cmdName, modeInvalid)
	}

	m := mode(os.Args[1])

	initModeToFlags(&mToFs, &fValues, cmdName)

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
		fs := mToFs[modeRender]
		fs.Parse(remainingArgs)

		if fs.NArg() > 1 {
			printUsageAndExit(cmdName, modeRender)
		}

		// Initialize renderer args from parsed values
		a := renderArgs{
			outDir:        fValues.render.outDir,
			force:         fValues.render.force,
			rendererCount: fValues.render.rendererCount,
		}

		// Initialize inDir to (optionally) first non-flag arg
		a.inDir = fs.Arg(0)
		if a.inDir == "" {
			a.inDir = "."
		}

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
		if len(remainingArgs) > 1 {
			printUsageAndExit(cmdName, modePreview)
		}
		inDir := filepath.Clean(remainingArgs[0])
		a := previewArgs{path: inDir}
		err := preview(a)
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

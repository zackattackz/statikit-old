package main

import (
	"fmt"
	"os"
)

type mode string

const (

	// Modes of operation
	modeInvalid mode = ""
	modeRender  mode = "render"
	modePreview mode = "preview"
	modePublish mode = "publish"
	modeInit    mode = "init"
	modeHelp    mode = "help"
)

func isMode(s string) (exists bool) {
	modeMap := map[mode]any{
		modeInvalid: nil,
		modeRender:  nil,
		modePreview: nil,
		modePublish: nil,
		modeInit:    nil,
		modeHelp:    nil,
	}
	_, exists = modeMap[mode(s)]
	return
}

func printUsageAndExit(cmdName string, m mode) {

	var opts string
	usageFmtStr := "usage: %s %s %s"
	fs := mToFs[m]
	allModes := fmt.Sprintf("[%s | %s | %s | %s | %s]", modeRender, modePreview, modePublish, modeInit, modeHelp)
	allModesButHelp := allModes[:len(allModes)-len(modeHelp)-4] + "]"

	switch m {
	case modeRender:
		opts = "[-o dirname] [-renderer-count uint] [-f] [dirname]"
	case modePreview:
		opts = "[dirname]"
	case modePublish:
		opts = "[dirname]"
	case modeInit:
		opts = "[dirname]"
	case modeHelp:
		opts = allModesButHelp
	case modeInvalid:
		m = mode(allModes)
	default:
		panic(fmt.Sprintf("invalid mode: %s", m))
	}
	usageStr := fmt.Sprintf(usageFmtStr, cmdName, m, opts)
	fmt.Fprintln(os.Stderr, usageStr)
	fs.PrintDefaults()
	os.Exit(2)
}

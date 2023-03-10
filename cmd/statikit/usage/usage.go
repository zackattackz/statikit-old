package usage

import (
	"fmt"
	"os"
	"path/filepath"
)

type Mode string

const (
	// Modes of operation
	Invalid Mode = ""
	Render  Mode = "render"
	Preview Mode = "preview"
	Init    Mode = "init"
	Help    Mode = "help"
)

func (m Mode) IsValid() bool {
	switch m {
	case Invalid, Render, Preview, Init, Help:
		return true
	default:
		return false
	}
}

func PrintUsageAndExit(cmdName string, m Mode, modeToPrintDefaults map[Mode]func()) {

	var opts string
	_, cmdName = filepath.Split(cmdName)
	usageFmtStr := "usage: %s %s %s"
	allModes := fmt.Sprintf("[%s | %s | %s | %s]", Render, Preview, Init, Help)
	allModesButHelp := allModes[:len(allModes)-len(Help)-4] + "]"

	switch m {
	case Render:
		opts = "[-o dirname] [-renderer-count uint] [-f] [dirname]"
	case Preview:
		opts = "[dirname]"
	case Init:
		opts = "[dirname]"
	case Help:
		opts = allModesButHelp
	case Invalid:
		m = Mode(allModes)
	default:
		panic(fmt.Sprintf("invalid mode: %s", m))
	}

	usageStr := fmt.Sprintf(usageFmtStr, cmdName, m, opts)
	fmt.Fprintln(os.Stderr, usageStr)
	if printDefaults, ok := modeToPrintDefaults[m]; ok {
		printDefaults()
	}
	os.Exit(2)
}

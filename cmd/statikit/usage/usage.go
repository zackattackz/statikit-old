package usage

import (
	"flag"
	"fmt"
	"os"
)

type Mode string

const (
	// Modes of operation
	Invalid Mode = ""
	Render  Mode = "render"
	Preview Mode = "preview"
	Publish Mode = "publish"
	Init    Mode = "init"
	Help    Mode = "help"
)

func IsMode(m Mode) (exists bool) {
	modeMap := map[Mode]any{
		Invalid: nil,
		Render:  nil,
		Preview: nil,
		Publish: nil,
		Init:    nil,
		Help:    nil,
	}
	_, exists = modeMap[m]
	return
}

func PrintUsageAndExit(cmdName string, m Mode, flagSet *flag.FlagSet) {

	var opts string
	usageFmtStr := "usage: %s %s %s"
	allModes := fmt.Sprintf("[%s | %s | %s | %s | %s]", Render, Preview, Publish, Init, Help)
	allModesButHelp := allModes[:len(allModes)-len(Help)-4] + "]"

	switch m {
	case Render:
		opts = "[-o dirname] [-renderer-count uint] [-f] [dirname]"
	case Preview:
		opts = "[dirname]"
	case Publish:
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
	if flagSet != nil {
		flagSet.PrintDefaults()
	}
	os.Exit(2)
}

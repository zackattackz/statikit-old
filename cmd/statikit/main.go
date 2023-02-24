package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/zackattackz/statikit-old/cmd/statikit/runners"
	"github.com/zackattackz/statikit-old/cmd/statikit/runners/help"
	"github.com/zackattackz/statikit-old/cmd/statikit/runners/initialize"
	"github.com/zackattackz/statikit-old/cmd/statikit/runners/preview"
	"github.com/zackattackz/statikit-old/cmd/statikit/runners/render"
	"github.com/zackattackz/statikit-old/cmd/statikit/usage"
	"github.com/zackattackz/statikit-old/internal/statikit/initializer"
	"github.com/zackattackz/statikit-old/internal/statikit/previewer"
	"github.com/zackattackz/statikit-old/internal/statikit/renderer"
)

type mainDependencies struct {
	modeToRunner        map[usage.Mode]runners.Runner
	modeToPrintDefaults map[usage.Mode]func()
}

func logErrAndExit(err error, code int) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(code)
}
func runMain(args []string, deps mainDependencies) {
	_, cmdName := filepath.Split(args[0])

	if len(args) < 2 {
		usage.PrintUsageAndExit(cmdName, usage.Invalid, nil)
	}

	m := usage.Mode(args[1])

	if !m.IsValid() {
		usage.PrintUsageAndExit(cmdName, usage.Invalid, nil)
	}

	usageForFunc := func(m usage.Mode) func() {
		return func() {
			usage.PrintUsageAndExit(cmdName, m, deps.modeToPrintDefaults)
		}
	}

	err := deps.modeToRunner[m](afero.NewOsFs(), args, usageForFunc)
	if err != nil {
		logErrAndExit(err, 1)
	} else {
		os.Exit(0)
	}
}

func main() {
	deps := mainDependencies{
		modeToRunner: map[usage.Mode]runners.Runner{
			usage.Preview: preview.Runner(previewer.Preview),
			usage.Render:  render.Runner(renderer.Render),
			usage.Init:    initialize.Runner(initializer.InitStatikitProject),
			usage.Help:    help.Runner,
		},
		modeToPrintDefaults: map[usage.Mode]func(){
			usage.Render: render.FlagSet.PrintDefaults,
		},
	}
	runMain(os.Args, deps)
}

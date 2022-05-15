package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zackattackz/azure_static_site_kit/pkg/statikit/config"
	"github.com/zackattackz/azure_static_site_kit/pkg/statikit/previewer"
	"github.com/zackattackz/azure_static_site_kit/pkg/statikit/renderer"
)

const (
	version     = "v0.1.0"
	usageString = "usage: statikit [render | preview | publish] opts"

	// Modes of operation
	render  = "render"
	preview = "preview"
	publish = "publish"

	// flag names
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

func logErrAndExit(err error, code int) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
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
	var rendererCount uint
	defaultOut := filepath.Join(os.TempDir(), "statikitRendered")
	flag.StringVar(&inDir, inFlag, defaultIn, descIn)
	flag.StringVar(&outDir, outFlag, defaultOut, descOut)
	flag.BoolVar(&force, forceFlag, defaultForce, descForce)
	flag.UintVar(&rendererCount, rendererCountFlag, defaultRendererCount, descRendererCount)
	flag.Parse()

	// If invalid mode print usage string and exit
	switch mode {
	case render, preview, publish:
		break
	default:
		fmt.Fprintln(os.Stderr, usageString)
		flag.PrintDefaults()
		os.Exit(2)
	}

	// Clean the in/out dirs
	inDir = filepath.Clean(inDir)
	outDir = filepath.Clean(outDir)

	// Ensure in dir exists and is a dir
	s, err := os.Stat(inDir)
	if err != nil {
		logErrAndExit(fmt.Errorf("couldn't read %s: %w", inDir, err), 1)
	}
	if !s.IsDir() {
		logErrAndExit(fmt.Errorf("%s is not a directory", inDir), 1)
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

	// Get the config path, open the file, and parse it
	configPath, configFormat, err := config.GetConfigPath(inDir)
	if err != nil {
		logErrAndExit(err, 1)
	}
	configFile, err := os.Open(configPath)
	if err != nil {
		logErrAndExit(err, 1)
	}
	config, err := config.ParseConfig(config.ParseConfigArgs{Reader: configFile, Format: configFormat})
	configFile.Close()
	if err != nil {
		logErrAndExit(err, 1)
	}

	// Determine cfg file name
	_, cfgFileName := filepath.Split(configPath)

	// Call the renderer
	rendererArgs := renderer.RendererArgs{
		InDir:         inDir,
		OutDir:        outDir,
		RendererCount: rendererCount,
		CfgFileName:   cfgFileName,
		Data:          config.Data,
	}
	if err := renderer.Render(rendererArgs); err != nil {
		logErrAndExit(err, 1)
	}

	// Go into mode specific handlers
	switch mode {

	case render:

	case preview:
		previewer.Preview(outDir)

	case publish:

	}

}

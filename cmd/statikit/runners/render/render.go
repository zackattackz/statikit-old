package render

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/runners"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/usage"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/config"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/initializer"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/renderer"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/schema"
)

var FlagSet *flag.FlagSet

var outDir, inDir string
var force bool
var rendererCount uint

// Initialize FlagSet
func init() {
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

	FlagSet = flag.NewFlagSet(string(usage.Render), flag.ExitOnError)
	defaultOut := filepath.Join(os.TempDir(), "statikitRendered")

	FlagSet.StringVar(&outDir, outFlag, defaultOut, descOut)
	FlagSet.BoolVar(&force, forceFlag, defaultForce, descForce)
	FlagSet.UintVar(&rendererCount, rendererCountFlag, defaultRendererCount, descRendererCount)

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

func Runner(render renderer.RenderFunc) runners.Runner {
	return func(args []string, usageFor runners.UsageForFunc) error {
		FlagSet.Usage = usageFor(usage.Render)

		FlagSet.Parse(args[2:])

		if FlagSet.NArg() > 1 {
			usageFor(usage.Render)()
		}

		// Initialize inDir to (optionally) first non-flag arg
		inDir = FlagSet.Arg(0)
		if inDir == "" {
			inDir = "."
		}

		// Clean the in/out dirs
		inDir = filepath.Clean(inDir)
		outDir = filepath.Clean(outDir)

		// Ensure in dir exists and is a dir
		inDirStat, err := os.Stat(inDir)
		if err != nil {
			return fmt.Errorf("couldn't read %s: %w", inDir, err)
		}
		if !inDirStat.IsDir() {
			return fmt.Errorf("%s is not a directory", inDir)
		}

		// If no force flag, ensure user wants to erase.
		if !force {
			if err := warnErase(outDir); err != nil {
				return err
			}
		}

		// If we make it here, erase outdir
		if err := os.RemoveAll(outDir); err != nil {
			return err
		}

		// Parse the schema map
		schemaMap := make(schema.Map)
		schemaParser := schema.NewParser(inDir)
		err = schemaParser.Parse(&schemaMap)
		if err != nil {
			return err
		}

		// Parse the config
		cfgParser, err := config.NewParser(inDir)
		if err != nil {
			return err
		}
		cfg := config.T{}
		cfgParser.Parse(&cfg)

		cfg.Ignore = append(cfg.Ignore, initializer.StatikitDirName)

		rendererArgs := renderer.Args{
			InDir:         inDir,
			OutDir:        outDir,
			RendererCount: rendererCount,
			SchemaMap:     schemaMap,
			Ignore:        cfg.Ignore,
		}
		return render(rendererArgs)
	}
}

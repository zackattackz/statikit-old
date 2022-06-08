package render

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/runners"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/usage"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/config"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/initializer"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/renderer"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/schema"
)

type renderFor string

const (
	renderLocal renderFor = renderFor("local")
	renderAzure renderFor = renderFor("az")
)

func (r renderFor) IsValid() bool {
	switch r {
	case renderLocal, renderAzure:
		return true
	default:
		return false
	}
}

var FlagSet *flag.FlagSet

var outDir, inDir string
var force bool
var rendererCount uint
var forString string

// Initialize FlagSet
func init() {
	const (
		// flags
		outFlag           = "o"
		forceFlag         = "f"
		rendererCountFlag = "renderer-count"
		forFlag           = "for"

		// default flag values
		defaultForce         = false
		defaultRendererCount = 20
		defaultFor           = renderLocal

		// flag descriptions
		descOut           = "rendered output directory"
		descForce         = "force output directory removal"
		descRendererCount = "how many renderer goroutines to be made"
		descFor           = "determines additional files to output"
	)

	FlagSet = flag.NewFlagSet(string(usage.Render), flag.ExitOnError)
	defaultOut := filepath.Join(os.TempDir(), "statikitRendered")

	FlagSet.StringVar(&outDir, outFlag, defaultOut, descOut)
	FlagSet.BoolVar(&force, forceFlag, defaultForce, descForce)
	FlagSet.UintVar(&rendererCount, rendererCountFlag, defaultRendererCount, descRendererCount)
	FlagSet.StringVar(&forString, forFlag, string(defaultFor), descFor)
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

type renderFunc func(renderer.Args) error

func Runner(render renderFunc) runners.Runner {
	return func(fs afero.Fs, args []string, usageFor runners.UsageForFunc) error {
		FlagSet.Usage = usageFor(usage.Render)

		FlagSet.Parse(args[2:])

		if FlagSet.NArg() > 1 {
			usageFor(usage.Render)()
		}

		rFor := renderFor(forString)
		if !rFor.IsValid() {
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
		inDirStat, err := fs.Stat(inDir)
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
		if err := fs.RemoveAll(outDir); err != nil {
			return err
		}

		// Parse the schema map
		schemaMap := make(schema.Map)
		schemaParser := schema.NewParser(fs, inDir)
		err = schemaParser.Parse(&schemaMap)
		if err != nil {
			return err
		}

		// Parse the config
		cfgParser, err := config.NewParser(fs, inDir)
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
			Fs:            fs,
		}
		err = render(rendererArgs)
		if err != nil {
			return err
		}

		switch rFor {
		case renderLocal:
		case renderAzure:
			// Copy over the keyfile, it is needed for when we publish
			outKeyFilePath := filepath.Join(outDir, initializer.StatikitDirName, initializer.KeyFileName)
			inKeyFilePath := filepath.Join(inDir, initializer.StatikitDirName, initializer.KeyFileName)
			err = fs.Mkdir(filepath.Join(outDir, initializer.StatikitDirName), 0755)
			if err != nil {
				return err
			}

			inF, err := fs.Open(inKeyFilePath)
			if err != nil {
				return err
			}
			defer inF.Close()

			outF, err := fs.Create(outKeyFilePath)
			if err != nil {
				return err
			}
			defer outF.Close()

			_, err = io.Copy(outF, inF)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

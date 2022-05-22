package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/zackattackz/azure_static_site_kit/internal/statikit/config"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/renderer"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/schema"
)

type renderArgs struct {
	inDir         string
	outDir        string
	force         bool
	rendererCount uint
	schemaMap     schema.Map
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

func render(a renderArgs) error {

	// Ensure in dir exists and is a dir
	s, err := os.Stat(a.inDir)
	if err != nil {
		return fmt.Errorf("couldn't read %s: %w", a.inDir, err)
	}
	if !s.IsDir() {
		return fmt.Errorf("%s is not a directory", a.inDir)
	}

	// If no force flag, ensure user wants to erase.
	if !a.force {
		if err := warnErase(a.outDir); err != nil {
			return err
		}
	}

	// If we make it here, erase outdir
	if err := os.RemoveAll(a.outDir); err != nil {
		return err
	}

	// Call the renderer
	rendererArgs := renderer.Args{
		InDir:         a.inDir,
		OutDir:        a.outDir,
		RendererCount: a.rendererCount,
		CfgDirName:    config.ConfigDirName,
		SchemaMap:     a.schemaMap,
	}
	return renderer.Run(rendererArgs)
}

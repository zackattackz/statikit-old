package main

import (
	"os"
	"path/filepath"

	"github.com/zackattackz/azure_static_site_kit/internal/statikit/config"
)

func parseConfig(inDir string) (*config.T, error) {
	// Get the config path, open the file, and parse it
	cfgPath, configFormat, err := config.GetPath(inDir)
	if err != nil {
		return nil, err
	}
	cfgFile, err := os.Open(filepath.Join(cfgPath))
	if err != nil {
		return nil, err
	}
	cfg, err := config.Parse(config.ParseArgs{Reader: cfgFile, Format: configFormat})
	cfgFile.Close()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

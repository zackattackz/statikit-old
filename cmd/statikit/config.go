package main

import (
	"os"
	"path/filepath"

	"github.com/zackattackz/azure_static_site_kit/pkg/statikit/config"
)

func parseConfig(inDir string) (*config.T, error) {
	// Get the config path, open the file, and parse it
	configPath, configFormat, err := config.GetPath(inDir)
	if err != nil {
		return nil, err
	}
	configFile, err := os.Open(filepath.Join(configPath))
	if err != nil {
		return nil, err
	}
	config, err := config.Parse(config.ParseArgs{Reader: configFile, Format: configFormat})
	configFile.Close()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

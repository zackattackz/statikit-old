package main

import (
	"os"

	"github.com/zackattackz/azure_static_site_kit/pkg/statikit/config"
)

func parseConfig(inDir string) (*config.T, string, error) {
	// Get the config path, open the file, and parse it
	configPath, configFormat, err := config.GetPath(inDir)
	if err != nil {
		return nil, "", err
	}
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, "", err
	}
	config, err := config.Parse(config.ParseArgs{Reader: configFile, Format: configFormat})
	configFile.Close()
	if err != nil {
		return nil, "", err
	}

	return &config, configPath, nil
}

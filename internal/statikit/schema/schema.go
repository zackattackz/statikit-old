package schema

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/config"
	sp "github.com/zackattackz/azure_static_site_kit/pkg/subtractPaths"
)

type Format uint

type T struct {
	Data map[string]any
}

// Maps path names to their data
type Map map[string]T

const (
	JsonFormat Format = iota
	TomlFormat

	DataDirName     = "schema"
	DefaultDataName = "_defaultvalues"
)

var extToFormat = map[string]Format{
	".json": JsonFormat,
	".toml": TomlFormat,
}

func parse(r io.Reader, format Format) (d T, err error) {
	switch format {
	case JsonFormat:
		dec := json.NewDecoder(r)
		err = dec.Decode(&d)
	case TomlFormat:
		dec := toml.NewDecoder(r)
		_, err = dec.Decode(&d)
	default:
		err = errors.New("invalid format")
	}
	if d.Data == nil {
		err = errors.New("parsed data is <nil>")
	}
	return
}

func Parse(root string) (Map, error) {
	res := make(Map)
	dataPath := filepath.Join(root, config.ConfigDirName, DataDirName)
	err := filepath.WalkDir(dataPath, func(path string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if e.IsDir() ||
			!e.Type().IsRegular() {
			return nil
		}

		pathFromData := sp.SubtractPaths(dataPath, path)
		ext := filepath.Ext(pathFromData)
		format := extToFormat[ext]
		pathWithoutExt := strings.TrimSuffix(pathFromData, ext)

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		d, err := parse(f, format)
		if err != nil {
			return err
		}
		res[pathWithoutExt] = d
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

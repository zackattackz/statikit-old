package schema

import (
	"encoding/json"
	"errors"
	"html/template"
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
	Data    map[string]any           // Variable names->raw data to be substituted, comes directly from schema
	FileSub map[string]template.HTML // Variable names->html to be substituted, comes from a file
}

type parseT struct {
	Data    map[string]any    // Variable names->raw data to be substituted
	FileSub map[string]string // Variable names->filename, relative to _statikit/.., that contains data to be substituted
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

func parse(r io.Reader, format Format) (d parseT, err error) {
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
		// Ensure there was no error in call
		if err != nil {
			return err
		}

		// Skip e if it is dir or non-regular
		if e.IsDir() ||
			!e.Type().IsRegular() {
			return nil
		}

		// Determine path without extension, to be used to address res
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

		// Populate a new schema.T with the parsed fields
		s := T{}
		s.Data = d.Data
		s.FileSub = make(map[string]template.HTML)

		// Read all the files in FileSubst and
		// fill out T's FileSubst with contents
		for v, fname := range d.FileSub {
			f, err := os.Open(filepath.Join(root, filepath.Clean(fname)))
			if err != nil {
				return err
			}
			b, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			s.FileSub[v] = template.HTML(b)
		}

		res[pathWithoutExt] = s
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type ConfigFileFormat uint

const (
	JsonFormat ConfigFileFormat = iota
	TomlFormat

	ConfigFileName = "statikitConfig"
)

var extToFormat = map[string]ConfigFileFormat{
	".json": JsonFormat,
	".toml": TomlFormat,
}

type NotExistError struct {
	path string
}

func (e NotExistError) Error() string {
	return fmt.Sprintf("config file does not exist at path %v", e.path)
}

func (e NotExistError) Is(target error) bool {
	targetCast, ok := target.(NotExistError)
	if !ok {
		return false
	}
	return targetCast.path == e.path
}

type MoreThanOneError struct {
	amount uint
}

func (e MoreThanOneError) Error() string {
	return fmt.Sprintf("too many config files: %v", e.amount)
}

func (e MoreThanOneError) Is(target error) bool {
	targetCast, ok := target.(MoreThanOneError)
	if !ok {
		return false
	}
	return targetCast.amount == e.amount
}

type T struct {
	Data any // The data passed to template renders
}

type ParseConfigArgs struct {
	Reader io.Reader
	Format ConfigFileFormat
}

/*
Searches for a single valid config file in path `root`

Returns:
	(string, ConfigFileFormat, nil) The full path to the single config file and its format
	(_, _, MoreThanOneErr) If more than one config file exists
	(_, _, ErrConfigFileNotExist) If no config files exist
	(_, _, error) Any generic error from os calls
*/
func GetConfigPath(root string) (resPath string, f ConfigFileFormat, resErr error) {
	// For each valid extention, ext,
	// If the file at path "root/p.ext" exists and is regular,
	// return the path and it's associated format
	// if no valid file is found return fs.ErrNotExist
	count := uint(0)
	for ext := range extToFormat {
		p := filepath.Join(root, ConfigFileName+ext)
		s, err := os.Stat(p)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		} else if err != nil {
			resErr = err
			return
		} else if !s.Mode().IsRegular() {
			continue
		} else {
			resPath = p
			f = extToFormat[ext]
			count += 1
		}
	}
	if count == 0 {
		resErr = NotExistError{path: root}
	} else if count > 1 {
		resErr = MoreThanOneError{amount: count}
	}
	return
}

func ParseConfig(a ParseConfigArgs) (result T, err error) {
	switch a.Format {
	case JsonFormat:
		d := json.NewDecoder(a.Reader)
		err = d.Decode(&result)
	case TomlFormat:
		d := toml.NewDecoder(a.Reader)
		_, err = d.Decode(&result)
	default:
		err = errors.New("invalid format")
	}
	if result.Data == nil {
		err = errors.New("parsed data is <nil>")
	}
	return
}

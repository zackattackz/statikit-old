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
	"github.com/zackattackz/azure_static_site_kit/pkg/secret"
)

type Format uint

const (
	JsonFormat Format = iota
	TomlFormat

	ConfigFileName = "config"
	ConfigDirName  = "_statikit"
	KeyFileName    = "key.aes256"
)

var extToFormat = map[string]Format{
	".json": JsonFormat,
	".toml": TomlFormat,
}

var formatToExt = map[Format]string{
	JsonFormat: ".json",
	TomlFormat: ".toml",
}

var formatToInit = map[Format]string{
	JsonFormat: "{\n\"Data\": {}\n}",
	TomlFormat: "[Data]",
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
	Ignore []string // List of file globs to ignore when rendering
}

type ParseArgs struct {
	Reader io.Reader
	Format Format
}

/*
Searches for a single valid config file in path `root`

Returns:
	(string, ConfigFileFormat, nil) The full path to the single config file and its format
	(_, _, MoreThanOneErr) If more than one config file exists
	(_, _, ErrConfigFileNotExist) If no config files exist
	(_, _, error) Any generic error from os calls
*/
func GetPath(root string) (resPath string, f Format, resErr error) {
	// For each valid extention, ext,
	// If the file at path "root/p.ext" exists and is regular,
	// return the path and it's associated format
	// if no valid file is found return fs.ErrNotExist
	count := uint(0)
	for ext := range extToFormat {
		p := filepath.Join(root, ConfigDirName, ConfigFileName+ext)
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

func Parse(a ParseArgs) (result T, err error) {
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
	// Clean all Ignore paths
	for i, p := range result.Ignore {
		result.Ignore[i] = filepath.Clean(p)
	}
	return
}

func Create(path string, f Format, pwd string, key []byte) error {

	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("%s already exists", path)
	}

	if err = os.Mkdir(path, 0755); err != nil {
		return err
	}

	path = filepath.Join(path, ConfigDirName)

	if err = os.Mkdir(path, 0755); err != nil {
		return err
	}

	cfgFile, err := os.Create(filepath.Join(path, ConfigFileName+formatToExt[f]))
	if err != nil {
		return err
	}
	defer cfgFile.Close()
	cfgFile.WriteString(formatToInit[f])

	keyFile, err := os.Create(filepath.Join(path, KeyFileName))
	if err != nil {
		return err
	}
	defer keyFile.Close()

	aes, err := secret.Encrypt(pwd, key)
	if err != nil {
		return err
	}

	n, err := keyFile.Write(aes)
	if err != nil {
		return err
	}
	for n < len(aes) {
		m, err := keyFile.Write(aes[n:])
		if err != nil {
			return err
		}
		n += m
	}
	return nil
}

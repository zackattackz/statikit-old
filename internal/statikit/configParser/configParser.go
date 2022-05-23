package configParser

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/initializer"
)

type Interface interface {
	Parse(*Config) error
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

type Config struct {
	Ignore []string // List of file globs to ignore when rendering
}

type t struct {
	root string
	path string
}

func New(root string) (Interface, error) {
	parser := &t{root: root}
	p, err := getPath(root)
	if err != nil {
		return nil, err
	}
	parser.path = p
	return parser, nil
}

/*
Searches for valid config file in path `root`
*/
func getPath(root string) (string, error) {
	p := filepath.Join(root, initializer.StatikitDirName, initializer.ConfigFileName)
	s, err := os.Stat(p)
	if errors.Is(err, fs.ErrNotExist) || !s.Mode().IsRegular() {
		return "", NotExistError{path: root}
	} else if err != nil {
		return "", err
	}
	return p, nil
}

func (t *t) Parse(cfg *Config) error {
	f, err := os.Open(t.path)
	if err != nil {
		return err
	}
	defer f.Close()
	d := toml.NewDecoder(f)
	_, err = d.Decode(cfg)
	if err != nil {
		return err
	}
	// Clean all Ignore paths
	for i, p := range cfg.Ignore {
		// If any pattern is bad, return ErrBadPattern
		if _, err = filepath.Match(p, ""); errors.Is(err, filepath.ErrBadPattern) {
			return err
		}
		cfg.Ignore[i] = filepath.Clean(p)
	}
	return nil
}

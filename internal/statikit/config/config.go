package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/afero"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/initializer"
)

type Parser interface {
	Parse(*T) error
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

type azblobConfig struct {
	AccountName   string
	ContainerName string
}

type T struct {
	Ignore []string     // List of file globs to ignore when rendering
	Az     azblobConfig // Details of az blob storage
}

type parser struct {
	root string
	path string
	fs   afero.Fs
}

func NewParser(fs afero.Fs, root string) (Parser, error) {
	parser := &parser{root: root, fs: fs}
	p, err := getPath(fs, root)
	if err != nil {
		return nil, err
	}
	parser.path = p
	return parser, nil
}

/*
Searches for valid config file in path `root`
*/
func getPath(fs afero.Fs, root string) (string, error) {
	p := filepath.Join(root, initializer.StatikitDirName, initializer.ConfigFileName)
	s, err := fs.Stat(p)
	if os.IsNotExist(err) || !s.Mode().IsRegular() {
		return "", NotExistError{path: root}
	} else if err != nil {
		return "", err
	}
	return p, nil
}

func (p *parser) Parse(cfg *T) error {
	f, err := p.fs.Open(p.path)
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

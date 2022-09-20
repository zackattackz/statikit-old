package schema

import (
	"errors"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/afero"
	"codeberg.org/zackattackz/statikit/internal/statikit/initializer"
	sp "codeberg.org/zackattackz/statikit/pkg/subtractPaths"
)

type Parser interface {
	Parse(*Map) error
}

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

func parse(r io.Reader) (d parseT, err error) {
	dec := toml.NewDecoder(r)
	_, err = dec.Decode(&d)
	if d.Data == nil {
		err = errors.New("parsed data is <nil>")
	}
	return
}

type parser struct {
	root string
	fs   afero.Fs
}

func NewParser(fs afero.Fs, root string) Parser {
	return &parser{root: root, fs: fs}
}

func (p *parser) Parse(m *Map) error {
	dataPath := filepath.Join(p.root, initializer.StatikitDirName, initializer.SchemaDirName)
	err := afero.Walk(p.fs, dataPath, func(path string, e fs.FileInfo, err error) error {
		// Ensure there was no error in call
		if err != nil {
			return err
		}

		// Determine path without extension, to be used to address res
		pathFromData := sp.SubtractPaths(dataPath, path)
		ext := filepath.Ext(pathFromData)
		pathWithoutExt := strings.TrimSuffix(pathFromData, ext)

		// Skip e if it is dir or non-regular or is not a .toml file
		if e.IsDir() ||
			!e.Mode().IsRegular() ||
			ext != ".toml" {
			return nil
		}

		f, err := p.fs.Open(path)
		if err != nil {
			return err
		}
		d, err := parse(f)
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
			f, err := p.fs.Open(filepath.Join(p.root, filepath.Clean(fname)))
			if err != nil {
				return err
			}
			b, err := afero.ReadAll(f)
			if err != nil {
				return err
			}
			s.FileSub[v] = template.HTML(b)
		}

		(*m)[pathWithoutExt] = s
		return nil
	})
	return err
}

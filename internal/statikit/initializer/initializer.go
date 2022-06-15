package initializer

import (
	"path/filepath"

	"github.com/spf13/afero"
)

const (
	StatikitDirName   = "_statikit"
	ConfigFileName    = "config.toml"
	SchemaDirName     = "schema"
	KeyFileName       = "key.aes256"
	DefaultValuesName = "_defaultvalues"
)

func InitStatikitProject(fs afero.Fs, path string) error {

	if err := fs.Mkdir(path, 0755); err != nil {
		return err
	}

	path = filepath.Join(path, StatikitDirName)

	if err := fs.Mkdir(path, 0755); err != nil {
		return err
	}

	cfgFile, err := fs.Create(filepath.Join(path, ConfigFileName))
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	return fs.Mkdir(filepath.Join(path, SchemaDirName), 0755)
}

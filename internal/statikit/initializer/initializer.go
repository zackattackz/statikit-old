package initializer

import (
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/zackattackz/azure_static_site_kit/pkg/secret"
)

const (
	StatikitDirName   = "_statikit"
	ConfigFileName    = "config.toml"
	SchemaDirName     = "schema"
	KeyFileName       = "key.aes256"
	DefaultValuesName = "_defaultvalues"
)

func InitStatikitProject(fs afero.Fs, path string, pwd string, key []byte) error {

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

	fs.Mkdir(filepath.Join(path, SchemaDirName), 0755)

	aes, err := secret.Encrypt(pwd, key)
	if err != nil {
		return err
	}

	return afero.WriteFile(fs, filepath.Join(path, KeyFileName), aes, 0755)
}

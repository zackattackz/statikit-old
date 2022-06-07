package initializer

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/zackattackz/azure_static_site_kit/pkg/secret"
)

const (
	StatikitDirName = "_statikit"
	ConfigFileName  = "config.toml"
	SchemaDirName   = "schema"
	KeyFileName     = "key.aes256"
	DefaultDataName = "_defaultvalues"
)

func InitStatikitProject(fs afero.Fs, path string, pwd string, key []byte) error {

	_, err := fs.Stat(path)
	if err == nil {
		return fmt.Errorf("%s already exists", path)
	}

	if err = fs.Mkdir(path, 0755); err != nil {
		return err
	}

	path = filepath.Join(path, StatikitDirName)

	if err = fs.Mkdir(path, 0755); err != nil {
		return err
	}

	cfgFile, err := fs.Create(filepath.Join(path, ConfigFileName))
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	keyFile, err := fs.Create(filepath.Join(path, KeyFileName))
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

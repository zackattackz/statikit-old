package initializer

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

func assertDirExists(t *testing.T, fs afero.Fs, path string) {
	info, err := fs.Stat(path)
	if err != nil {
		t.Fatalf("error on fs.Stat(%s): %v", path, err)
	}
	if !info.IsDir() {
		t.Fatalf("%s is not a directory", path)
	}
}

func TestInitStatikitProject(t *testing.T) {
	fs := afero.NewMemMapFs()

	testPath := "test"

	err := InitStatikitProject(fs, testPath)
	if err != nil {
		t.Fatalf("error on InitStatikitProject(): %v", err)
	}

	assertDirExists(t, fs, testPath)
	assertDirExists(t, fs, filepath.Join(testPath, StatikitDirName))
	assertDirExists(t, fs, filepath.Join(testPath, StatikitDirName, SchemaDirName))
}

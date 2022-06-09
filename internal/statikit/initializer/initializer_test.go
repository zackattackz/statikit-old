package initializer

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/zackattackz/azure_static_site_kit/pkg/secret"
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

	type testInput struct {
		pwd string
		key []byte
	}

	testInputs := []testInput{
		{
			pwd: "testPassword",
			key: []byte("testKey"),
		},
	}

	expectedAESs := make([][]byte, len(testInputs))

	for i, testInput := range testInputs {
		expectedAES, err := secret.Encrypt(testInput.pwd, testInput.key)
		if err != nil {
			t.Fatalf("error on secret.Encrypt(%v, %v): %v", testInput.pwd, testInput.key, err)
		}
		expectedAESs[i] = expectedAES
	}

	for i, input := range testInputs {
		testPath := fmt.Sprint(i)
		err := InitStatikitProject(fs, testPath, input.pwd, input.key)
		if err != nil {
			t.Fatalf("error on InitStatikitProject(): %v", err)
		}

		assertDirExists(t, fs, testPath)
		assertDirExists(t, fs, filepath.Join(testPath, StatikitDirName))
		assertDirExists(t, fs, filepath.Join(testPath, StatikitDirName, SchemaDirName))

		keyFilePath := filepath.Join(testPath, StatikitDirName, KeyFileName)
		f, err := fs.Open(keyFilePath)
		if err != nil {
			t.Fatalf("error on Open(%s): %v", keyFilePath, err)
		}

		actualAES, err := afero.ReadAll(f)
		if err != nil {
			t.Fatalf("error on ReadAll for file %s: %v", keyFilePath, err)
		}

		expectedKey, err := secret.Decrypt(input.pwd, expectedAESs[i])
		if err != nil {
			t.Fatalf("error on secret.Decrypt(%v): %v", expectedAESs[i], err)
		}
		actualKey, err := secret.Decrypt(input.pwd, actualAES)
		if err != nil {
			t.Fatalf("error on secret.Decrypt(%v): %v", actualAES, err)
		}
		if !bytes.Equal(expectedKey, actualKey) {
			t.Fatalf("expected key not equal to actual: %v ..... %v", expectedKey, actualKey)
		}

	}

}

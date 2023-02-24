package config

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/afero"
	"github.com/zackattackz/statikit-old/internal/statikit/initializer"
)

func testParse(t *testing.T, testInputs []string, expectedResults []T) {
	fs := afero.NewMemMapFs()

	for i, testInput := range testInputs {
		fpath := filepath.Join(fmt.Sprint(i), initializer.StatikitDirName, initializer.ConfigFileName)
		dname, _ := filepath.Split(fpath)
		err := fs.MkdirAll(dname, 0755)
		if err != nil {
			t.Fatalf("error creating test input directory: %s, %v", dname, err)
		}
		err = afero.WriteFile(fs, fpath, []byte(testInput), 0755)
		if err != nil {
			t.Fatalf("error creating test input file: %s, %v", fpath, err)
		}
	}

	for i, expectedResult := range expectedResults {
		testName := fmt.Sprint(i)
		cfgParser, err := NewParser(fs, testName)
		if err != nil {
			t.Fatalf("error on New(%s): %v", testName, err)
		}
		cfg := T{}
		err = cfgParser.Parse(&cfg)
		if err != nil {
			t.Fatalf("error on Parse(): %v", err)
		}
		if !reflect.DeepEqual(cfg, expectedResult) {
			t.Fatalf("expected %v, actual %v", expectedResult, cfg)
		}
	}
}

func TestParse(t *testing.T) {
	// Test Ignore
	func() {
		testInputs := []string{
			"",
			"Ignore = []",
			"Ignore = [ \"templates/\" ]",
			"Ignore = [ \"someDir/templates\" ]",
			"Ignore = [ \"one\", \"two/three\", \"four\" ]",
		}

		expectedResults := []T{
			{},
			{Ignore: []string{}},
			{Ignore: []string{filepath.Clean("templates/")}},
			{Ignore: []string{filepath.Join("someDir", "templates")}},
			{Ignore: []string{"one", filepath.Join("two", "three"), "four"}},
		}
		testParse(t, testInputs, expectedResults)
	}()

	// Test Az
	func() {
		testInputs := []string{
			"[Az]\nAccountName = \"Test\"",
			"[Az]\nContainerName = \"Test\"",
			"[Az]\nAccountName = \"Test\"\nContainerName = \"Test\"",
		}

		expectedResults := []T{
			{Az: azblobConfig{AccountName: "Test"}},
			{Az: azblobConfig{ContainerName: "Test"}},
			{Az: azblobConfig{AccountName: "Test", ContainerName: "Test"}},
		}
		testParse(t, testInputs, expectedResults)
	}()
}

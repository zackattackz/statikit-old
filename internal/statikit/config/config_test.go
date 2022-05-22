package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func testParseIgnore(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error on os.Getwd(): %s", err)
	}
	testcasesPath := filepath.Join(wd, "testcases")
	in := filepath.Join(testcasesPath, "in")

	expectedResults := map[string][]string{
		"six": {filepath.Clean("templates/")},
	}

	for testName, expectedResult := range expectedResults {
		testPath := filepath.Join(in, testName)
		d, err := os.Stat(testPath)
		if err != nil {
			t.Fatalf("error on Stat(%s): %v", testPath, err)
		}
		if !d.IsDir() {
			t.Fatalf("entry is not a directory: %s", d.Name())
		}
		p, f, err := GetPath(testPath)
		if err != nil {
			t.Fatalf("error on GetPath(%s): %v", testPath, err)
		}
		cfgFile, err := os.Open(p)
		if err != nil {
			t.Fatalf("error on Open(%s): %v", p, err)
		}
		actual, err := Parse(ParseArgs{Reader: cfgFile, Format: f})
		if err != nil {
			t.Fatalf("error on Parse(): %v", err)
		}
		if !reflect.DeepEqual(actual.Ignore, expectedResult) {
			t.Fatalf("expected %v, actual %v", expectedResult, actual.Ignore)
		}
	}
}

func TestParse(t *testing.T) {
	testParseIgnore(t)
}

type expectedResults struct {
	p   string
	f   Format
	err error
}

func TestGetPath(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error on os.Getwd(): %s", err)
	}
	testcasesPath := filepath.Join(wd, "testcases")
	in := filepath.Join(testcasesPath, "in")

	expectedResults := map[string]expectedResults{
		"one": {
			p:   filepath.Join(in, "one", ConfigDirName, ConfigFileName+".json"),
			f:   JsonFormat,
			err: nil,
		},
		"two": {
			p:   filepath.Join(in, "two", ConfigDirName, ConfigFileName+".toml"),
			f:   TomlFormat,
			err: nil,
		},
		"three": {
			p:   "",
			f:   0,
			err: NotExistError{path: filepath.Join(in, "three")},
		},
		"four": {
			p:   "",
			f:   0,
			err: NotExistError{path: filepath.Join(in, "four")},
		},
		"five": {
			p:   "",
			f:   0,
			err: MoreThanOneError{amount: 2},
		},
	}

	for testName, expectedResult := range expectedResults {
		testPath := filepath.Join(in, testName)
		d, err := os.Stat(testPath)
		if err != nil {
			t.Fatalf("error on Stat(%s): %v", testPath, err)
		}
		if !d.IsDir() {
			t.Fatalf("entry is not a directory: %s", d.Name())
		}

		p, f, err := GetPath(testPath)
		if err != nil {
			if !errors.Is(expectedResult.err, err) {
				t.Fatalf("expected.err = %v, actual err = %v", expectedResult.err, err)
			}
		} else {
			if expectedResult.p != p {
				t.Fatalf("expected.p = %v, actual p = %v", expectedResult.p, p)
			} else if expectedResult.f != f {
				t.Fatalf("expected.f = %v, actual f = %v", expectedResult.f, f)
			}
		}
	}
}

package statikit

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type runTestArgs struct {
	input    string
	expected map[string]interface{}
}

func runTest(a runTestArgs, format ParseDataFormat) error {
	r := strings.NewReader(a.input)
	actual, err := ParseDataFile(ParseDataArgs{Reader: r, Format: format})
	if err != nil {
		return err
	}
	actualMap := actual.(map[string]interface{})
	if !reflect.DeepEqual(actualMap, a.expected) {
		return fmt.Errorf("expected: \"%v\", actual: \"%v\"", a.expected, actualMap)
	}
	return nil
}

func TestParseData(t *testing.T) {
	tomlTests := []runTestArgs{
		{
			input:    `Test = "hello"`,
			expected: map[string]interface{}{"Test": "hello"},
		},
		{
			input:    "One = 1\nTwo = 2",
			expected: map[string]interface{}{"One": int64(1), "Two": int64(2)},
		},
	}

	jsonTests := []runTestArgs{
		{
			input:    `{"Test": "hello"}`,
			expected: map[string]interface{}{"Test": "hello"},
		},
		{
			input:    `{"One": 1, "Two": 2}`,
			expected: map[string]interface{}{"One": float64(1), "Two": float64(2)},
		},
	}

	for _, jsonTest := range jsonTests {
		err := runTest(jsonTest, JsonFormat)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, tomlTest := range tomlTests {
		err := runTest(tomlTest, TomlFormat)
		if err != nil {
			t.Fatal(err)
		}
	}
}

type expectedResults struct {
	p   string
	f   ParseDataFormat
	err error
}

func TestGetDataFilePath(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error on os.Getwd(): %s", err)
	}
	testcasesPath := filepath.Join(wd, "testcases", "datafile")
	in := filepath.Join(testcasesPath, "in")

	expectedResults := map[string]expectedResults{
		"one": {
			p:   filepath.Join(in, "one", DataFileName+".json"),
			f:   JsonFormat,
			err: nil,
		},
		"two": {
			p:   filepath.Join(in, "two", DataFileName+".toml"),
			f:   TomlFormat,
			err: nil,
		},
		"three": {
			p:   "",
			f:   0,
			err: ErrDataFileNotExist,
		},
		"four": {
			p:   "",
			f:   0,
			err: ErrDataFileNotExist,
		},
		"five": {
			p:   "",
			f:   0,
			err: moreThanOneError(2),
		},
	}

	d, err := os.ReadDir(in)
	if err != nil {
		t.Fatalf("error on ReadDir(\"%s\"): %s", in, err)
	}

	for _, e := range d {
		if !e.IsDir() {
			t.Fatalf("entry is not a directory: %s", e.Name())
		}

		in := filepath.Join(in, e.Name())

		p, f, err := GetDataFilePath(in)
		expected := expectedResults[e.Name()]
		if err != nil {
			if !errors.Is(expected.err, err) {
				t.Fatalf("expected.err = %v, actual err = %v", expected.err, err)
			}
		} else {
			if expected.p != p {
				t.Fatalf("expected.p = %v, actual p = %v", expected.p, p)
			} else if expected.f != f {
				t.Fatalf("expected.f = %v, actual f = %v", expected.f, f)
			}
		}
	}
}

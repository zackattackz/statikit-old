package statikit

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func dirsEqual(a, b string) (bool, error) {
	var item1, item2 string
	notEqualErr := errors.New("contents not equal")

	err := filepath.WalkDir(a, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		path = subtractPaths(a, path)
		fullA := filepath.Join(a, path)
		fullB := filepath.Join(b, path)

		if d.IsDir() {
			// Assert that fullB is a dir too
			fB, err := os.Open(fullB)
			if err != nil {
				return err
			}
			fBStat, err := fB.Stat()
			if err != nil {
				return err
			}
			if fBStat.IsDir() {
				return nil
			} else {
				item1 = fullA
				item2 = fullB
				return notEqualErr
			}
		}

		// If file isn't regular, just skip it
		if !d.Type().IsRegular() {
			return nil
		}

		// Assert that files at fullA and fullB are the same
		fullAContents, err := os.ReadFile(fullA)
		if err != nil {
			return err
		}
		fullBContents, err := os.ReadFile(fullB)
		if err != nil {
			return err
		}

		if bytes.Equal(fullAContents, fullBContents) {
			return nil
		} else {
			item1 = fullA
			item2 = fullB
			return notEqualErr
		}

	})
	if err == notEqualErr {
		return false, fmt.Errorf("%s, %s: %w", item1, item2, err)
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func TestRender(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error on os.Getwd(): %s", err)
	}
	testcasesPath := filepath.Join(wd, "testcases", "renderer")
	fmt.Println(testcasesPath)
	in := filepath.Join(testcasesPath, "in")
	out := filepath.Join(testcasesPath, "out")
	expected := filepath.Join(testcasesPath, "expected")

	d, err := os.ReadDir(in)
	if err != nil {
		t.Fatalf("error on ReadDir(\"%s\"): %s", in, err)
	}

	for _, e := range d {
		if !e.IsDir() {
			t.Fatalf("entry is not a directory: %s", e.Name())
		}
		in := filepath.Join(in, e.Name())
		out := filepath.Join(out, e.Name())
		expected := filepath.Join(expected, e.Name())

		os.RemoveAll(out)
		// TODO: Implement Data reader
		args := RendererArgs{InDir: in, OutDir: out, RendererCount: 20, Data: struct{ TestThree string }{TestThree: "world"}}
		err := Render(args)
		if err != nil {
			t.Fatalf("error on Render(%v): %s", args, err)
		}
		areEqual, err := dirsEqual(out, expected)
		if err != nil {
			t.Fatalf("error on dirsEqual(\"%s\", \"%s\"): %s", out, expected, err)
		}
		if !areEqual {
			t.Fatalf("expected dir %s did not equal actual dir %s", expected, out)
		}
	}
}
package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/initializer"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/schema"
	sp "github.com/zackattackz/azure_static_site_kit/pkg/subtractPaths"
)

type notEqualErr struct {
	aName       string
	bName       string
	aContents   []byte
	bContents   []byte
	notSameType bool
}

func (e notEqualErr) Error() string {
	if e.notSameType {
		return fmt.Sprintf("%s not equal to %s: not same type", e.aName, e.bName)
	}
	return fmt.Sprintf("%s not equal to %s: not same contents\n%s contents:\n%s\n%s contents:\n%s\n", e.aName, e.bName, e.aName, e.aContents, e.bName, e.bContents)
}

func dirsEqual(aferoFs afero.Fs, a, b string) (bool, error) {
	err := afero.Walk(aferoFs, a, func(path string, e fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		path = sp.SubtractPaths(a, path)
		fullA := filepath.Join(a, path)
		fullB := filepath.Join(b, path)

		if e.IsDir() {
			// Assert that fullB is a dir too
			fB, err := aferoFs.Open(fullB)
			if err != nil {
				return err
			}
			defer fB.Close()
			fBStat, err := fB.Stat()
			if err != nil {
				return err
			}
			if fBStat.IsDir() {
				return nil
			} else {
				return notEqualErr{aName: fullA, bName: fullB, notSameType: true}
			}
		}

		// If file isn't regular, just skip it
		if !e.Mode().IsRegular() {
			return nil
		}

		// Assert that files at fullA and fullB are the same
		fullAContents, err := afero.ReadFile(aferoFs, fullA)
		if err != nil {
			return err
		}
		fullBContents, err := afero.ReadFile(aferoFs, fullB)
		if err != nil {
			return err
		}

		if bytes.Equal(fullAContents, fullBContents) {
			return nil
		} else {
			return notEqualErr{aName: fullA, bName: fullB, aContents: fullAContents, bContents: fullBContents, notSameType: false}
		}
	})
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func TestRun(t *testing.T) {
	type pathAndContents struct {
		path     string
		contents string
	}
	type testInput struct {
		schemaMap schema.Map
		ignore    []string
		files     []pathAndContents
	}
	fs := afero.NewMemMapFs()
	in := "in"
	expected := "expected"
	out := "out"

	testInputs := []testInput{
		{
			schema.Map{},
			[]string{},
			[]pathAndContents{{"hello.txt", "Hello, world!"}},
		},
		{
			schema.Map{},
			[]string{},
			[]pathAndContents{
				{filepath.Join(initializer.StatikitDirName, "config.toml"), ""},
				{filepath.Join(initializer.StatikitDirName, initializer.SchemaDirName, "ignore.toml"), "[Data]"},
			},
		},
		{
			schema.Map{},
			[]string{},
			[]pathAndContents{
				{filepath.Join("dir", "file"), ""},
			},
		},
		{
			schema.Map{},
			[]string{},
			[]pathAndContents{
				{filepath.Join("dir", "file1"), "1"},
				{filepath.Join("dir", "another", "file2"), "2"},
				{filepath.Join("dir", "another", "deep", "file3"), "3"},
				{filepath.Join("dir", "anotherTwin", "file4"), "4"},
				{filepath.Join("dir", "file5"), "5"},
				{filepath.Join("something", "else", "file6"), "6"},
				{"file7", "7"},
			},
		},
		{
			schema.Map{
				"hello": {
					Data: map[string]any{
						"Test": "world",
					},
					FileSub: map[string]template.HTML{},
				},
			},
			[]string{},
			[]pathAndContents{
				{"hello.gohtml", "Hello, {{.Data.Test}}!"},
			},
		},
		{
			schema.Map{
				initializer.DefaultValuesName: {
					Data: map[string]any{
						"Count": int64(123),
					},
					FileSub: map[string]template.HTML{},
				},
			},
			[]string{},
			[]pathAndContents{
				{"hello.gohtml", "One twenty three: {{.Data.Count}}"},
			},
		},
		{
			schema.Map{
				initializer.DefaultValuesName: {
					Data: map[string]any{
						"Count": int64(200),
					},
					FileSub: map[string]template.HTML{},
				},
				"blog": {
					Data: map[string]any{
						"Count": int64(100),
					},
					FileSub: map[string]template.HTML{},
				},
			},
			[]string{},
			[]pathAndContents{
				{"blog.gohtml", "One hundred: {{.Data.Count}}"},
			},
		},
		{
			schema.Map{
				initializer.DefaultValuesName: {
					Data: map[string]any{
						"Title": "Hello",
						"Date":  "Today",
					},
					FileSub: map[string]template.HTML{
						"Head": template.HTML("<head><title>Hello World</title></head>"),
					},
				},
			},
			[]string{"templates"},
			[]pathAndContents{
				{"blog.one.gohtml", "{{.FileSub.Head}}\nPost title: {{.Data.Title}} Post date: {{.Data.Date}}"},
				{filepath.Join("templates", "head.html"), "<head><title>Hello World</title></head>"},
			},
		},
	}

	expectedOutputs := [][]pathAndContents{
		{{"hello.txt", "Hello, world!"}},
		{},
		{{filepath.Join("dir", "file"), ""}},
		{
			{filepath.Join("dir", "file1"), "1"},
			{filepath.Join("dir", "another", "file2"), "2"},
			{filepath.Join("dir", "another", "deep", "file3"), "3"},
			{filepath.Join("dir", "anotherTwin", "file4"), "4"},
			{filepath.Join("dir", "file5"), "5"},
			{filepath.Join("something", "else", "file6"), "6"},
			{"file7", "7"},
		},
		{{"hello.html", "Hello, world!"}},
		{{"hello.html", "One twenty three: 123"}},
		{{"blog.html", "One hundred: 100"}},
		{{"blog.one.html", "<head><title>Hello World</title></head>\nPost title: Hello Post date: Today"}},
	}

	for i, testInput := range testInputs {
		for _, pathAndContents := range testInput.files {
			path := filepath.Join(in, fmt.Sprint(i), pathAndContents.path)
			dname, _ := filepath.Split(path)
			err := fs.MkdirAll(dname, 0755)
			if err != nil {
				t.Fatalf("error creating test input directory: %s, %v", dname, err)
			}
			err = afero.WriteFile(fs, path, []byte(pathAndContents.contents), 0755)
			if err != nil {
				t.Fatalf("error creating test input file: %s, %v", path, err)
			}
		}
	}

	for i, files := range expectedOutputs {
		expectedDirPath := filepath.Join(expected, fmt.Sprint(i))
		err := fs.MkdirAll(expectedDirPath, 0755)
		if err != nil {
			t.Fatalf("error creating test input directory: %s, %v", expectedDirPath, err)
		}
		for _, pathAndContents := range files {
			path := filepath.Join(expected, fmt.Sprint(i), pathAndContents.path)
			dname, _ := filepath.Split(path)
			err := fs.MkdirAll(dname, 0755)
			if err != nil {
				t.Fatalf("error creating test input directory: %s, %v", dname, err)
			}
			err = afero.WriteFile(fs, path, []byte(pathAndContents.contents), 0755)
			if err != nil {
				t.Fatalf("error creating test input file: %s, %v", path, err)
			}
		}
	}

	for i, testInput := range testInputs {
		inPath := filepath.Join(in, fmt.Sprint(i))
		outPath := filepath.Join(out, fmt.Sprint(i))
		expectedPath := filepath.Join(expected, fmt.Sprint(i))

		args := Args{
			InDir:         inPath,
			OutDir:        outPath,
			RendererCount: 20,
			SchemaMap:     testInput.schemaMap,
			Ignore:        testInput.ignore,
			Fs:            fs,
		}
		args.Ignore = append(args.Ignore, initializer.StatikitDirName)

		err := Render(args)
		if err != nil {
			t.Fatalf("error on Render(%v): %s", args, err)
		}

		areEqual, err := dirsEqual(fs, outPath, expectedPath)
		if err != nil {
			t.Fatalf("error on dirsEqual(\"%s\", \"%s\"): %s", out, expected, err)
		}
		if !areEqual {
			t.Fatalf("expected dir %s did not equal actual dir %s", expected, out)
		}
	}
}

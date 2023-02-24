package schema

import (
	"fmt"
	"html/template"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/afero"
	"github.com/zackattackz/statikit-old/internal/statikit/initializer"
)

func TestParse(t *testing.T) {
	type pathAndContents struct {
		path     string
		ext      string
		contents string
	}

	testInputs := [][]pathAndContents{
		{
			{"blog", "toml", "[Data]\nTestOne = \"TestOne\"\nTestTwo = 2"},
		},
		{
			{"blog", "toml", "[Data]\nTestOne = 100"},
			{filepath.Join("folder", "nested"), "toml", "[Data]\nNestedTest = \"Nested\""},
		},
		{
			{"blog", "toml", "[Data]\nBlogStuff = \"Hi There\""},
			{filepath.Join("folder", "nested"), "toml", "[Data]\nNestedTest = \"Nested\""},
			{filepath.Join("folder", "blog"), "toml", "[Data]\nJsonBlog = \"HelloWorld\"\nTitle = \"Title\"\nDate = \"Today\""},
			{filepath.Join("folder", "deep", "status"), "toml", "[Data]\nLikes = 100\nPerson = \"Jesus\""},
			{filepath.Join("another", "another.testing"), "toml", "[Data]\nSomeMore = \"Stuff\""},
		},
		{
			{"blog", "toml", "[Data]\nTestOne = \"TestOne\"\nTestTwo = 2\n[FileSub]\nHead = \"templates/head.html\""},
			{filepath.Join("templates", "head"), "html", "<head><title>Hello World</title></head>"},
		},
		{
			{initializer.DefaultValuesName, "toml", "[Data]\nTitle = \"One\""},
		},
	}

	expectedResults := []Map{

		{
			"blog": T{
				Data: map[string]any{
					"TestOne": "TestOne",
					"TestTwo": int64(2),
				},
				FileSub: map[string]template.HTML{},
			},
		},
		{
			"blog": T{
				Data: map[string]any{
					"TestOne": int64(100),
				},
				FileSub: map[string]template.HTML{},
			},
			filepath.Join("folder", "nested"): T{
				Data: map[string]any{
					"NestedTest": "Nested",
				},
				FileSub: map[string]template.HTML{},
			},
		},
		{
			"blog": T{
				Data: map[string]any{
					"BlogStuff": "Hi There",
				},
				FileSub: map[string]template.HTML{},
			},
			filepath.Join("another", "another.testing"): T{
				Data: map[string]any{
					"SomeMore": "Stuff",
				},
				FileSub: map[string]template.HTML{},
			},
			filepath.Join("folder", "nested"): T{
				Data: map[string]any{
					"NestedTest": "Nested",
				},
				FileSub: map[string]template.HTML{},
			},
			filepath.Join("folder", "blog"): T{
				Data: map[string]any{"JsonBlog": "HelloWorld",
					"Title": "Title",
					"Date":  "Today",
				},
				FileSub: map[string]template.HTML{},
			},
			filepath.Join("folder", "deep", "status"): T{
				Data: map[string]any{"Likes": int64(100),
					"Person": "Jesus",
				},
				FileSub: map[string]template.HTML{},
			},
		},
		{
			"blog": T{
				Data: map[string]any{
					"TestOne": "TestOne",
					"TestTwo": int64(2),
				},
				FileSub: map[string]template.HTML{
					"Head": "<head><title>Hello World</title></head>",
				},
			},
		},
		{
			initializer.DefaultValuesName: T{
				Data: map[string]any{
					"Title": "One",
				},
				FileSub: map[string]template.HTML{},
			},
		},
	}

	fs := afero.NewMemMapFs()

	for i, testInput := range testInputs {

		for _, f := range testInput {
			var fpath string
			if f.ext == "toml" {
				fpath = filepath.Join(fmt.Sprint(i), initializer.StatikitDirName, initializer.SchemaDirName, f.path+"."+f.ext)
			} else {
				fpath = filepath.Join(fmt.Sprint(i), f.path+"."+f.ext)
			}
			dname, _ := filepath.Split(fpath)
			err := fs.MkdirAll(dname, 0755)
			if err != nil {
				t.Fatalf("error creating test input directory: %s, %v", dname, err)
			}
			err = afero.WriteFile(fs, fpath, []byte(f.contents), 0755)
			if err != nil {
				t.Fatalf("error creating test input file: %s, %v", fpath, err)
			}
		}
	}

	for i, expectedResult := range expectedResults {
		testName := fmt.Sprint(i)
		schemaParser := NewParser(fs, testName)
		actual := make(Map)
		err := schemaParser.Parse(&actual)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(actual, expectedResult) {
			t.Fatalf("expected %v, actual %v", expectedResult, actual)
		}
	}
}

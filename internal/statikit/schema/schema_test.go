package schema

import (
	"html/template"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	expectedMap := map[string]Map{
		"one": {
			"blog": T{
				Data: map[string]any{
					"TestOne": "TestOne",
					"TestTwo": int64(2),
				},
				FileSub: map[string]template.HTML{},
			},
		},
		"two": {
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
		"three": {
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
		"four": {
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
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error on os.Getwd(): %s", err)
	}
	testcasesPath := filepath.Join(wd, "testcases")
	in := filepath.Join(testcasesPath, "in")

	e, err := os.ReadDir(in)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range e {
		in := filepath.Join(in, e.Name())
		schemaParser := NewParser(in)
		actual := make(Map)
		err := schemaParser.Parse(&actual)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(actual, expectedMap[e.Name()]) {
			t.Fatalf("expected %v, actual %v", expectedMap[e.Name()], actual)
		}
	}
}

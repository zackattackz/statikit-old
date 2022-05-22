package schema

import (
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
				FileSub: map[string]string{},
			},
		},
		"two": {
			"blog": T{
				Data: map[string]any{
					"TestOne": int64(100),
				},
				FileSub: map[string]string{},
			},
			filepath.Join("folder", "nested"): T{
				Data: map[string]any{
					"NestedTest": "Nested",
				},
				FileSub: map[string]string{},
			},
		},
		"three": {
			"blog": T{
				Data: map[string]any{
					"BlogStuff": "Hi There",
				},
				FileSub: map[string]string{},
			},
			filepath.Join("another", "another.testing"): T{
				Data: map[string]any{
					"SomeMore": "Stuff",
				},
				FileSub: map[string]string{},
			},
			filepath.Join("folder", "nested"): T{
				Data: map[string]any{
					"NestedTest": "Nested",
				},
				FileSub: map[string]string{},
			},
			filepath.Join("folder", "blog"): T{
				Data: map[string]any{"JsonBlog": "HelloWorld",
					"Title": "Title",
					"Date":  "Today",
				},
				FileSub: map[string]string{},
			},
			filepath.Join("folder", "deep", "status"): T{
				Data: map[string]any{"Likes": float64(100),
					"Person": "Jesus",
				},
				FileSub: map[string]string{},
			},
		},
		"four": {
			"blog": T{
				Data: map[string]any{
					"TestOne": "TestOne",
					"TestTwo": int64(2),
				},
				FileSub: map[string]string{
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
		actual, err := Parse(in)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(actual, expectedMap[e.Name()]) {
			t.Fatalf("expected %v, actual %v", expectedMap[e.Name()], actual)
		}
	}
}

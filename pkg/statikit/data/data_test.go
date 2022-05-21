package data

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	expectedMap := map[string]Map{
		"one": {
			"blog": {
				"TestOne": "TestOne",
				"TestTwo": int64(2),
			},
		},
		"two": {
			"blog": {
				"TestOne": int64(100),
			},
			filepath.Join("folder", "nested"): {
				"NestedTest": "Nested",
			},
		},
		"three": {
			"blog": {
				"BlogStuff": "Hi There",
			},
			filepath.Join("another", "another.testing"): {
				"SomeMore": "Stuff",
			},
			filepath.Join("folder", "nested"): {
				"NestedTest": "Nested",
			},
			filepath.Join("folder", "blog"): {
				"JsonBlog": "HelloWorld",
				"Title":    "Title",
				"Date":     "Today",
			},
			filepath.Join("folder", "deep", "status"): {
				"Likes":  float64(100),
				"Person": "Jesus",
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

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
			},
		},
		"two": {
			"blog": T{
				Data: map[string]any{
					"TestOne": int64(100),
				},
			},
			filepath.Join("folder", "nested"): T{
				Data: map[string]any{
					"NestedTest": "Nested",
				},
			},
		},
		"three": {
			"blog": T{
				Data: map[string]any{
					"BlogStuff": "Hi There",
				},
			},
			filepath.Join("another", "another.testing"): T{
				Data: map[string]any{
					"SomeMore": "Stuff",
				},
			},
			filepath.Join("folder", "nested"): T{
				Data: map[string]any{
					"NestedTest": "Nested",
				},
			},
			filepath.Join("folder", "blog"): T{
				Data: map[string]any{"JsonBlog": "HelloWorld",
					"Title": "Title",
					"Date":  "Today",
				},
			},
			filepath.Join("folder", "deep", "status"): T{
				Data: map[string]any{"Likes": float64(100),
					"Person": "Jesus",
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

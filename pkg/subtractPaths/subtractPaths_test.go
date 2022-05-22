package subtractpaths

import (
	"path/filepath"
	"testing"
)

func TestSubtractPaths(t *testing.T) {
	expected := []struct{ parent, child, expected string }{
		{
			parent:   filepath.Join("usr", "lib", "somelib"),
			child:    filepath.Join("usr", "lib", "somelib", "someFolder"),
			expected: filepath.Clean("someFolder"),
		},
		{
			parent:   filepath.Join("usr", "lib", "somelib"),
			child:    filepath.Join("usr", "lib", "somelib"),
			expected: "", // Specifically "", NOT "." - see filepath.Join doc
		},
		{
			parent:   filepath.Join("usr", "lib", "somelib"),
			child:    filepath.Join("usr", "lib", "somelib", "someFolder", "someFile.png"),
			expected: filepath.Join("someFolder", "someFile.png"),
		},
	}

	for _, testCase := range expected {
		actual := SubtractPaths(testCase.parent, testCase.child)
		if testCase.expected != actual {
			t.Fatalf("expected %v, actual %v", testCase.expected, actual)
		}
	}
}

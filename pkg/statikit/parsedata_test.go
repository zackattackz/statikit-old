package statikit

import (
	"fmt"
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
	actual, err := ParseData(ParseDataArgs{r: r, format: format})
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

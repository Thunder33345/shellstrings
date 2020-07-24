package shellstrings

import (
	"strings"
	"testing"
)

func TestMidwaySlash(t *testing.T) {
	input := `"Foo bar" hello \World`
	output := []string{"Foo bar", `hello`, `\World`}
	p := Parse(input)
	if !isEqual(p, output) {
		errorHelper(t, input, output, p)
	}
}

func TestTrailingBothQS(t *testing.T) {
	input := `\bar"\`
	output := []string{`\bar"\`}
	p := Parse(input)
	if !isEqual(p, output) {
		errorHelper(t, input, output, p)
	}
}

func TestTrailingBothSQ(t *testing.T) {
	input := `"Foo bar" hello World\"`
	output := []string{"Foo bar", `hello`, `World"`}
	p := Parse(input)
	if !isEqual(p, output) {
		errorHelper(t, input, output, p)
	}
}

func TestTrailingBothSSQ(t *testing.T) {
	input := `"Foo bar" hello World\\"`
	output := []string{"Foo bar", `hello`, `World\\"`}
	p := Parse(input)
	if !isEqual(p, output) {
		errorHelper(t, input, output, p)
	}
}

func TestTrailingBothQQS(t *testing.T) {
	input := `"Foo bar" hello World""\`
	output := []string{"Foo bar", `hello`, `World\`}
	p := Parse(input)
	if !isEqual(p, output) {
		errorHelper(t, input, output, p)
	}
}

func TestWeirdQuotes(t *testing.T) {
	tests := []TestTemplate{
		{"nothing in middle", `hel""lo`, []string{`hello`}},
		{"cross space", `hel" "lo`, []string{`hel lo`}},
		{"space cross space", `hel " " lo`, []string{`hel`, ` `, `lo`}},
		{"quotes n text", `h""e l" "" "l  o"""`, []string{`he`, `l  l`, `o"`}},
		{"just quotes", `"" " "" "  """`, []string{`  `, `"`}},
		{"midway quotes", `"Foo bar" "hello World`, []string{"Foo bar", `"hello`, `World`}},
	}
	for _, test := range tests {
		t.Run(t.Name(), testHelper(test))
	}
}

func TestEscapedQuotes(t *testing.T) {
	tests := []TestTemplate{
		{"simple escape quotes", `\"hi\"`, []string{`"hi"`}},
		{"escape quotes", `"hi hi\" lol`, []string{`"hi`, `hi"`, `lol`}},
		{"order quotes", `"hi ok\" k \"lol 12\"`, []string{`"hi`, `ok"`, `k`, `"lol`, `12"`}},
		{"leading quotes", `"hello world`, []string{`"hello`, `world`}},
		//Return: hello lol" haha
	}
	for _, test := range tests {
		t.Run(t.Name(), testHelper(test))
	}
}

func TestBasic(t *testing.T) {
	tests := []TestTemplate{
		{"simple", "hello", []string{"hello"}},
		{"quoted", `"hello world"`, []string{"hello world"}},
		{"multi", `"hello" "world"`, []string{"hello", "world"}},
		{"complicated", `"Foo bar" esca\"ped normal "to\\ns of s\\la\\sh"`, []string{"Foo bar", "esca\"ped", "normal", `to\\ns of s\\la\\sh`}},
		{"escape irrelevant", `"\F\oo bar" \s\o\m\e`, []string{`\F\oo bar`, `\s\o\m\e`}},
		{"escape quotes", `"hello lol\" haha`, []string{`"hello`, `lol"`, `haha`}},
		{"trailing slash", `hello World\`, []string{`hello`, `World\`}},
		{"trailing quote", `Quote"`, []string{`Quote"`}},
	}
	for _, test := range tests {
		t.Run(t.Name(), testHelper(test))
	}
}
func testHelper(test TestTemplate) func(*testing.T) {
	return func(t *testing.T) {
		p := Parse(test.input)
		if !isEqual(p, test.output) {
			t.Errorf("Test#%s\nInput: [%s]\nReturn: %v\nWant:   %v", test.name, test.input, strings.Join(p, ","), strings.Join(test.output, ","))
		}
	}
}

func errorHelper(t *testing.T, input string, expected []string, result []string) {
	t.Errorf("Input: Parse(%s)\nReturn: %v\nWant: %v", input, strings.Join(result, ","), strings.Join(expected, ","))
}

type TestTemplate struct {
	name   string
	input  string
	output []string
}

func isEqual(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

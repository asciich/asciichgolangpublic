package strings

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringsGetFirstLine(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"", ""},
		{"testcase", "testcase"},
		{"testcase\n", "testcase"},
		{"testcase\nanother line", "testcase"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				firstLine := GetFirstLine(tt.input)
				assert.EqualValues(tt.expectedOutput, firstLine)
			},
		)
	}
}

func TestStringsGetFirstLineAndTrimSpace(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"", ""},
		{"testcase", "testcase"},
		{"testcase\n", "testcase"},
		{"testcase\nanother line", "testcase"},
		{"\n", ""},
		{" \n", ""},
		{"\t\n", ""},
		{"\t \n", ""},
		{" testcase", "testcase"},
		{" testcase\n", "testcase"},
		{" testcase\nanother line", "testcase"},
		{"\ttestcase", "testcase"},
		{"\ttestcase\n", "testcase"},
		{"\ttestcase\nanother line", "testcase"},
		{"testcase ", "testcase"},
		{"testcase \n", "testcase"},
		{"testcase \nanother line", "testcase"},
		{"testcase\t", "testcase"},
		{"testcase\t\n", "testcase"},
		{"testcase\t\nanother line", "testcase"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				firstLine := GetFirstLineAndTrimSpace(tt.input)
				assert.EqualValues(tt.expectedOutput, firstLine)
			},
		)
	}
}

func TestStringsEnsureEndsWithExactlyOneLine(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"", "\n"},
		{"\n", "\n"},
		{"\n\n", "\n"},
		{"a", "a\n"},
		{"a\n", "a\n"},
		{"a\n\n", "a\n"},
		{"a\n\n\n", "a\n"},
		{"a\nb", "a\nb\n"},
		{"a\nb\n", "a\nb\n"},
		{"a\nb\n\n", "a\nb\n"},
		{"a\nb\n\n\n", "a\nb\n"},
		{"a\nb\n\n\n\n", "a\nb\n"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				ensuredLineBreak := EnsureEndsWithExactlyOneLineBreak(tt.input)
				assert.EqualValues(tt.expectedOutput, ensuredLineBreak)
			},
		)
	}
}

func TestStringsRemoveTailingNewline(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"", ""},
		{"\n", ""},
		{"a", "a"},
		{"a\n", "a"},
		{"ab\n", "ab"},
		{"abc\n", "abc"},
		{"ab", "ab"},
		{"abc", "abc"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				ensuredLineBreak := RemoveTailingNewline(tt.input)
				assert.EqualValues(tt.expectedOutput, ensuredLineBreak)
			},
		)
	}
}

func TestStringsEnsureFirstCharUppercase(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"", ""},
		{"a", "A"},
		{"A", "A"},
		{"abc", "Abc"},
		{"Abc", "Abc"},
		{"AbC", "AbC"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				firstCharUppercased := EnsureFirstCharUppercase(tt.input)
				assert.EqualValues(tt.expectedOutput, firstCharUppercased)
			},
		)
	}
}

func TestStringsEnsureFirstCharLowercase(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"", ""},
		{"a", "a"},
		{"A", "a"},
		{"abc", "abc"},
		{"Abc", "abc"},
		{"AbC", "abC"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				firstCharUppercased := EnsureFirstCharLowercase(tt.input)
				assert.EqualValues(tt.expectedOutput, firstCharUppercased)
			},
		)
	}
}

func TestStringsRemoveComments(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"", ""},
		{"abc", "abc"},
		{"abc\n", "abc\n"},
		{"abc\ndef", "abc\ndef"},
		{"abc\ndef\n", "abc\ndef\n"},
		{"abc\n#def\n", "abc\n"},
		{"#abc\n#def\n", ""},
		{"#abc\ndef\n", "def\n"},
		{"abc\n//def\n", "abc\n"},
		{"//abc\n//def\n", ""},
		{"//abc\ndef\n", "def\n"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				commentsRemoved := RemoveComments(tt.input)
				assert.EqualValues(tt.expectedOutput, commentsRemoved)
			},
		)
	}
}

func TestStringsRightFillWithSpaces(t *testing.T) {
	tests := []struct {
		input          string
		fillLenght     int
		expectedOutput string
	}{
		{"", 0, ""},
		{"", -1, ""},
		{"", -100, ""},
		{"", 6, "      "},
		{"a", 6, "a     "},
		{"ab", 6, "ab    "},
		{"abc", 6, "abc   "},
		{"abcd", 6, "abcd  "},
		{"abcde", 6, "abcde "},
		{"abcdef", 6, "abcdef"},
		{"abcdefg", 6, "abcdefg"},
		{"abcdefgh", 6, "abcdefgh"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				filled := RightFillWithSpaces(tt.input, tt.fillLenght)
				assert.EqualValues(tt.expectedOutput, filled)
			},
		)
	}
}

func TestStringsHasPrefixIgnoreCase(t *testing.T) {
	tests := []struct {
		input             string
		prefix            string
		expectedHasPrefix bool
	}{
		{"", "", true},
		{"abc", "a", true},
		{"abc", "A", true},
		{"abc", "Ab", true},
		{"abc", "aB", true},
		{"abc", "b", false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedHasPrefix,
					HasPrefixIgnoreCase(tt.input, tt.prefix),
				)
			},
		)
	}
}

func TestStringsTrimPrefixIgnoreCase(t *testing.T) {
	tests := []struct {
		input           string
		prefix          string
		expectedTrimmed string
	}{
		{"", "", ""},
		{"abc", "a", "bc"},
		{"abc", "A", "bc"},
		{"abc", "Ab", "c"},
		{"abc", "aB", "c"},
		{"abc", "b", "abc"},
		{"abc", "abc", ""},
		{"abc", "ABC", ""},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedTrimmed,
					TrimPrefixIgnoreCase(tt.input, tt.prefix),
				)
			},
		)
	}
}

func TestStringsIsFirstCharLowerCase(t *testing.T) {
	tests := []struct {
		input                      string
		expectedFirstCharLowerCase bool
	}{
		{"", false},
		{"abc", true},
		{"aBC", true},
		{"ABC", false},
		{"Abc", false},
		{" abc", false},
		{"_abc", false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedFirstCharLowerCase,
					IsFirstCharLowerCase(tt.input),
				)
			},
		)
	}
}

func TestStringsIsFirstCharUpperCase(t *testing.T) {
	tests := []struct {
		input                      string
		expectedFirstCharLowerCase bool
	}{
		{"", false},
		{"abc", false},
		{"aBC", false},
		{"ABC", true},
		{"Abc", true},
		{" abc", false},
		{"_abc", false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedFirstCharLowerCase,
					IsFirstCharUpperCase(tt.input),
				)
			},
		)
	}
}

func TestStringsSplitLines(t *testing.T) {
	tests := []struct {
		input         string
		expectedLines []string
	}{
		{"", []string{}},
		{"\n", []string{}},
		{"\n\n", []string{""}},
		{"hello", []string{"hello"}},
		{"hello\nworld", []string{"hello", "world"}},
		{"hello\r\nworld", []string{"hello", "world"}},
		{"hello\nworld\n", []string{"hello", "world"}},
		{"hello\nworld\n\n", []string{"hello", "world", ""}},
		{"hello\nworld\n\n\n", []string{"hello", "world", "", ""}},
		{"hello\nworld\n\nabc", []string{"hello", "world", "", "abc"}},
		{"hello\r\nworld\r\n", []string{"hello", "world"}},
		{"hello\r\nworld\r\n\r\n", []string{"hello", "world", ""}},
		{"hello\nworld\nworld2", []string{"hello", "world", "world2"}},
		{"hello\r\nworld\r\nworld2", []string{"hello", "world", "world2"}},
		{"hello\nworld\nworld2\n", []string{"hello", "world", "world2"}},
		{"hello\r\nworld\r\nworld2\r\n", []string{"hello", "world", "world2"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedLines,
					SplitLines(tt.input, true),
				)
			},
		)
	}
}


func TestStringsSplitWords(t *testing.T) {
	tests := []struct {
		input         string
		expectedWords []string
	}{
		{"", []string{}},
		{" ", []string{}},
		{"hello", []string{"hello"}},
		{"hello world", []string{"hello", "world"}},
		{"hello (world){}", []string{"hello", "world"}},
		{"hello (.world){}", []string{"hello", "world"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedWords,
					SplitWords(tt.input),
				)
			},
		)
	}
}

func TestStrings_MatchesRegex(t *testing.T) {
	tests := []struct {
		input         string
		regex         string
		expectedMatch bool
	}{
		{"abc", "abc", true},
		{"abc", "^abc", true},
		{"abc", "^abc$", true},
		{"abc", "^abcd$", false},
		{"a.log", "a.log", true},
		{"ablog", "a.log", true},
		{"ablog", "a\\.log", false},
		{"a.log", ".*.log", true},
		{"a.log", ".*\\.log", true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedMatch,
					MustMatchesRegex(tt.input, tt.regex),
				)
			},
		)
	}
}

func TestStringsIsComment(t *testing.T) {
	tests := []struct {
		input             string
		expectedIsComment bool
	}{
		{"", false},
		{" ", false},
		{"hello", false},
		{"hello world", false},
		{"#hello world", true},
		{"# hello world", true},
		{"# hello world\n", true},
		{"# REPLACE_BETWEEN_MARKERS START source=./stages.txt", true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedIsComment,
					IsComment(tt.input),
				)
			},
		)
	}
}

func TestStringsTrimSpacesLeft(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"", ""},
		{" ", ""},
		{" a", "a"},
		{" abc", "abc"},
		{"\ta", "a"},
		{"\tabc", "abc"},
		{"\na", "a"},
		{"\nabc", "abc"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedOutput,
					TrimSpacesLeft(tt.input),
				)
			},
		)
	}
}

func TestStringsContainsAtLeastOneSubstring(t *testing.T) {
	tests := []struct {
		input            string
		subsrings        []string
		expectedContains bool
	}{
		{"", []string{}, false},
		{"a", []string{"a"}, true},
		{"a", []string{"a", "b"}, true},
		{"a", []string{"z", "a", "b"}, true},
		{"A", []string{"a"}, false},
		{"A", []string{"a", "b"}, false},
		{"A", []string{"z", "a", "b"}, false},
		{"ABC", []string{"a"}, false},
		{"ABC", []string{"a", "b"}, false},
		{"ABC", []string{"z", "a", "b"}, false},
		{"aBC", []string{"a"}, true},
		{"aBC", []string{"a", "b"}, true},
		{"aBC", []string{"z", "a", "b"}, true},
		{"iJc", []string{"a"}, false},
		{"iJc", []string{"a", "b"}, false},
		{"iJc", []string{"z", "a", "b"}, false},
		{"IJC", []string{"a"}, false},
		{"IJC", []string{"a", "b"}, false},
		{"IJC", []string{"z", "a", "b"}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedContains,
					ContainsAtLeastOneSubstring(tt.input, tt.subsrings),
				)
			},
		)
	}
}

func TestContainsAtLeastOneSubstringIngoreCase(t *testing.T) {
	tests := []struct {
		input            string
		subsrings        []string
		expectedContains bool
	}{
		{"", []string{}, false},
		{"a", []string{"a"}, true},
		{"a", []string{"a", "b"}, true},
		{"a", []string{"z", "a", "b"}, true},
		{"A", []string{"a"}, true},
		{"A", []string{"a", "b"}, true},
		{"A", []string{"z", "a", "b"}, true},
		{"ABC", []string{"a"}, true},
		{"ABC", []string{"a", "b"}, true},
		{"ABC", []string{"z", "a", "b"}, true},
		{"aBC", []string{"a"}, true},
		{"aBC", []string{"a", "b"}, true},
		{"aBC", []string{"z", "a", "b"}, true},
		{"iJc", []string{"a"}, false},
		{"iJc", []string{"a", "b"}, false},
		{"iJc", []string{"z", "a", "b"}, false},
		{"IJC", []string{"a"}, false},
		{"IJC", []string{"a", "b"}, false},
		{"IJC", []string{"z", "a", "b"}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedContains,
					ContainsAtLeastOneSubstringIgnoreCase(tt.input, tt.subsrings),
				)
			},
		)
	}
}

func TestStringsContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		input            string
		subsring         string
		expectedContains bool
	}{
		{"hello WORLD", "hallo", false},
		{"hello WORLD", "HALLO", false},
		{"hello WORLD", "hello", true},
		{"hello WORLD", "HELLO", true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedContains,
					ContainsIgnoreCase(tt.input, tt.subsring),
				)
			},
		)
	}
}

func TestStringsTrimAllLeadingAndTailingNewLines(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"", ""},
		{"testcase", "testcase"},
		{"testcase\n", "testcase"},
		{"\ntestcase", "testcase"},
		{"\ntestcase\n", "testcase"},
		{"\ntestcase\n\n", "testcase"},
		{"\n\ntestcase\n\n", "testcase"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				output := TrimAllLeadingAndTailingNewLines(tt.input)
				assert.EqualValues(tt.expectedOutput, output)
			},
		)
	}
}

func TestStrings_RemoveLinesWithPrefix(t *testing.T) {
	tests := []struct {
		input          string
		prefix         string
		expectedOutput string
	}{
		{"", "abc", ""},
		{"\n", "abc", "\n"},
		{"abc\n", "abc", ""},
		{"1: a\n2: b\n3: c\n", "1", "2: b\n3: c\n"},
		{"1: a\n2: b\n3: c", "1", "2: b\n3: c"},
		{"1: a\n2: b\n3: c\n", "2", "1: a\n3: c\n"},
		{"1: a\n2: b\n3: c", "2", "1: a\n3: c"},
		{"1: a\n2: b\n3: c\n", "2:", "1: a\n3: c\n"},
		{"1: a\n2: b\n3: c", "2:", "1: a\n3: c"},
		{"1: a\n2: b\n3: c\n", "2: ", "1: a\n3: c\n"},
		{"1: a\n2: b\n3: c", "2: ", "1: a\n3: c"},
		{"1: a\n2: b\n3: c\n", "3", "1: a\n2: b\n"},
		{"1: a\n2: b\n3: c", "3", "1: a\n2: b"},
		{"1: a\n2: b\n3: c\n", "3:", "1: a\n2: b\n"},
		{"1: a\n2: b\n3: c", "3:", "1: a\n2: b"},
		{"1: a\n2: b\n3: c\n", "3: ", "1: a\n2: b\n"},
		{"1: a\n2: b\n3: c", "3: ", "1: a\n2: b"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedOutput,
					RemoveLinesWithPrefix(
						tt.input,
						tt.prefix,
					),
				)
			},
		)
	}
}

func TestStrings_HexStringToBytes(t *testing.T) {
	tests := []struct {
		hexString string
		hexBytes  []byte
	}{
		{"", []byte{}},
		{"0", []byte{0}},
		{"00", []byte{0}},
		{"0x00", []byte{0}},
		{"0X00", []byte{0}},
		{"1", []byte{1}},
		{"01", []byte{1}},
		{"0x01", []byte{1}},
		{"0X01", []byte{1}},
		{"a", []byte{10}},
		{"0a", []byte{10}},
		{"0x0a", []byte{10}},
		{"0X0a", []byte{10}},
		{"A", []byte{10}},
		{"0A", []byte{10}},
		{"0x0A", []byte{10}},
		{"0X0A", []byte{10}},
		{"0a00", []byte{10, 0}},
		{"0x0a00", []byte{10, 0}},
		{"0X0a00", []byte{10, 0}},
		{"0A00", []byte{10, 0}},
		{"0x0A00", []byte{10, 0}},
		{"0X0A00", []byte{10, 0}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.hexBytes,
					MustHexStringToBytes(tt.hexString),
				)
			},
		)
	}
}

func TestStrings_ContainsLine(t *testing.T) {
	tests := []struct {
		input            string
		line             string
		expectedContains bool
	}{
		{"", "", false},
		{"a\nb", "", false},
		{"a\n\nb", "", true},
		{"a\nb\nc", "a", true},
		{"a\nb\nc", "b", true},
		{"a\nb\nc", "c", true},
		{"a\nb\nc", "bc", false},
		{"a\nhello world\nc", "hello world", true},
		{"a\nhello world\nc", "hello world ", false},
		{"a\nhello world\nc", " hello world ", false},
		{"a\nhello world\nc", "hello", false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedContains,
					ContainsLine(tt.input, tt.line),
				)
			},
		)
	}
}

func TestStrings_GetAsKeyValues(t *testing.T) {
	tests := []struct {
		input             string
		expectedKeyValues map[string]string
	}{
		{"", map[string]string{}},
		{"\n", map[string]string{}},
		{"a=b", map[string]string{"a": "b"}},
		{"a=b\n", map[string]string{"a": "b"}},
		{"a:b", map[string]string{"a": "b"}},
		{"a:b\n", map[string]string{"a": "b"}},
		{" a=b", map[string]string{"a": "b"}},
		{" a=b\n", map[string]string{"a": "b"}},
		{" a:b", map[string]string{"a": "b"}},
		{" a:b\n", map[string]string{"a": "b"}},
		{"a =b", map[string]string{"a": "b"}},
		{"a =b\n", map[string]string{"a": "b"}},
		{"a :b", map[string]string{"a": "b"}},
		{"a :b\n", map[string]string{"a": "b"}},
		{"a = b", map[string]string{"a": "b"}},
		{"a = b\n", map[string]string{"a": "b"}},
		{"a : b", map[string]string{"a": "b"}},
		{"a : b\n", map[string]string{"a": "b"}},
		{"a = b ", map[string]string{"a": "b"}},
		{"a = b \n", map[string]string{"a": "b"}},
		{"a : b ", map[string]string{"a": "b"}},
		{"a : b \n", map[string]string{"a": "b"}},
		{"\na=b", map[string]string{"a": "b"}},
		{"\na=b\nc=d", map[string]string{"a": "b", "c": "d"}},
		{"\na:b", map[string]string{"a": "b"}},
		{"\na:b\nc:d", map[string]string{"a": "b", "c": "d"}},
		{"\na=b\n", map[string]string{"a": "b"}},
		{"\na=b\nc=d\n", map[string]string{"a": "b", "c": "d"}},
		{"\na:b\n", map[string]string{"a": "b"}},
		{"\na:b\nc:d\n", map[string]string{"a": "b", "c": "d"}},
		{"\na=b\n\n\n\n", map[string]string{"a": "b"}},
		{"\na=b\n\n\n\nc=d\n", map[string]string{"a": "b", "c": "d"}},
		{"\na:b\n\n\n\n", map[string]string{"a": "b"}},
		{"\na:b\n\n\n\nc:d\n", map[string]string{"a": "b", "c": "d"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedKeyValues,
					MustGetAsKeyValues(tt.input),
				)
			},
		)
	}
}

func TestStrings_GetValueAsString(t *testing.T) {
	tests := []struct {
		input         string
		key           string
		expectedValue string
	}{
		{"a=b\nc=hello world\n", "a", "b"},
		{"a=b\nc=hello world\n", "c", "hello world"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedValue,
					MustGetValueAsString(tt.input, tt.key),
				)
			},
		)
	}
}

func TestStrings_GetValueAsInt(t *testing.T) {
	tests := []struct {
		input         string
		key           string
		expectedValue int
	}{
		{"a=15\nb=0\nc=-3\n", "a", 15},
		{"a=15\nb=0\nc=-3\n", "b", 0},
		{"a=15\nb=0\nc=-3\n", "c", -3},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedValue,
					MustGetValueAsInt(tt.input, tt.key),
				)
			},
		)
	}
}

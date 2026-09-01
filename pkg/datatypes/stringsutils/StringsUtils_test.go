package stringsutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
				firstLine := GetFirstLine(tt.input)
				require.EqualValues(t, tt.expectedOutput, firstLine)
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
				firstLine := GetFirstLineAndTrimSpace(tt.input)
				require.EqualValues(t, tt.expectedOutput, firstLine)
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
				ensuredLineBreak := EnsureEndsWithExactlyOneLineBreak(tt.input)
				require.EqualValues(t, tt.expectedOutput, ensuredLineBreak)
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
				ensuredLineBreak := RemoveTailingNewline(tt.input)
				require.EqualValues(t, tt.expectedOutput, ensuredLineBreak)
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
				firstCharUppercased := EnsureFirstCharUppercase(tt.input)
				require.EqualValues(t, tt.expectedOutput, firstCharUppercased)
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
				firstCharUppercased := EnsureFirstCharLowercase(tt.input)
				require.EqualValues(t, tt.expectedOutput, firstCharUppercased)
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
				commentsRemoved := RemoveComments(tt.input)
				require.EqualValues(t, tt.expectedOutput, commentsRemoved)
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
				filled := RightFillWithSpaces(tt.input, tt.fillLenght)
				require.EqualValues(t, tt.expectedOutput, filled)
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
				require.EqualValues(t, tt.expectedHasPrefix, HasPrefixIgnoreCase(tt.input, tt.prefix))
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
				require.EqualValues(t, tt.expectedTrimmed, TrimPrefixIgnoreCase(tt.input, tt.prefix))
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
				require.EqualValues(t, tt.expectedFirstCharLowerCase, IsFirstCharLowerCase(tt.input))
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
				require.EqualValues(t, tt.expectedFirstCharLowerCase, IsFirstCharUpperCase(tt.input))
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
		{"\n", []string{""}},
		{"\n\n", []string{"", ""}},
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
				require.EqualValues(t, tt.expectedLines, SplitLines(tt.input, true))
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
				require.EqualValues(t, tt.expectedWords, SplitWords(tt.input))
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
				require.EqualValues(t, tt.expectedMatch, MustMatchesRegex(tt.input, tt.regex))
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
				require.EqualValues(t, tt.expectedIsComment, IsComment(tt.input))
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
				require.EqualValues(t, tt.expectedOutput, TrimSpacesLeft(tt.input))
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
				require.EqualValues(t, tt.expectedContains, ContainsAtLeastOneSubstring(tt.input, tt.subsrings))
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
				require.EqualValues(t, tt.expectedContains, ContainsAtLeastOneSubstringIgnoreCase(tt.input, tt.subsrings))
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
				require.EqualValues(t, tt.expectedContains, ContainsIgnoreCase(tt.input, tt.subsring))
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
				output := TrimAllLeadingAndTailingNewLines(tt.input)
				require.EqualValues(t, tt.expectedOutput, output)
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
				require.EqualValues(t, tt.expectedOutput, RemoveLinesWithPrefix(tt.input, tt.prefix))
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
				require.EqualValues(t, tt.hexBytes, MustHexStringToBytes(tt.hexString))
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
				require.EqualValues(t, tt.expectedContains, ContainsLine(tt.input, tt.line))
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
				require.EqualValues(t, tt.expectedKeyValues, MustGetAsKeyValues(tt.input))
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
				require.EqualValues(t, tt.expectedValue, MustGetValueAsString(tt.input, tt.key))
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
				require.EqualValues(t, tt.expectedValue, MustGetValueAsInt(tt.input, tt.key))
			},
		)
	}
}

func TestStrings_EnsureSuffix(t *testing.T) {
	tests := []struct {
		input    string
		suffix   string
		expected string
	}{
		{"a", "\n", "a\n"},
		{"a\n", "\n", "a\n"},
		{"a", "x", "ax"},
		{"a\n", "x", "a\nx"},
		{"ax", "x", "ax"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expected, EnsureSuffix(tt.input, tt.suffix))
			},
		)
	}
}

func TestStrings_EnsurePrefix(t *testing.T) {
	tests := []struct {
		input    string
		prefix   string
		expected string
	}{
		{"a", "\n", "\na"},
		{"a\n", "\n", "\na\n"},
		{"\na", "\n", "\na"},
		{"a", "x", "xa"},
		{"a\n", "x", "xa\n"},
		{"xa", "x", "xa"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expected, EnsurePrefix(tt.input, tt.prefix))
			},
		)
	}
}

func TestStrings_ToHexString(t *testing.T) {
	tests := []struct {
		input     string
		delimiter string
		expected  string
	}{
		{"", "", ""},
		{"ls", "", "6c73"},
		{"ls", " ", "6c 73"},
		{"ls", "_", "6c_73"},
		{"ls", "_-_", "6c_-_73"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expected, ToHexString(tt.input, tt.delimiter))
			},
		)
	}
}

func TestStrings_ToHexStringSlice(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", []string{}},
		{"ls", []string{"6c", "73"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expected, ToHexStringSlice(tt.input))
			},
		)
	}
}

func Test_IsBeforeInAlphabeth(t *testing.T) {
	tests := []struct {
		input    string
		input2   string
		expected bool
	}{
		{"", "", false},
		{"a", "", false},
		{"a", "a", false},
		{"a", "b", true},
		{"aa", "b", true},
		{"aa", "bb", true},
		{"b", "a", false},
		{"b", "aa", false},
		{"bb", "aa", false},
		{"1", "b", true},
		{"1", "1", false},
		{"11", "b", true},
		{"11", "bb", true},
		{"b", "1", false},
		{"b", "11", false},
		{"bb", "11", false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expected, IsBeforeInAlphabeth(tt.input, tt.input2))
			},
		)
	}
}

func Test_AddLinePrefix(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := AddLinePrefix("", "")
		require.EqualValues(t, "", got)
	})

	t.Run("empty content", func(t *testing.T) {
		got := AddLinePrefix("", "prefix")
		require.EqualValues(t, "prefix", got)
	})

	t.Run("empty prefix", func(t *testing.T) {
		got := AddLinePrefix("hello", "")
		require.EqualValues(t, "hello", got)
	})

	t.Run("empty prefix with tailing newline for content", func(t *testing.T) {
		got := AddLinePrefix("hello\n", "")
		require.EqualValues(t, "hello\n", got)
	})

	t.Run("no newline at end", func(t *testing.T) {
		input := `abc
def`
		expected := `    abc
    def`

		got := AddLinePrefix(input, "    ")
		require.EqualValues(t, expected, got)
	})

	t.Run("newline at end", func(t *testing.T) {
		input := `abc
def
`
		expected := `    abc
    def
`
		got := AddLinePrefix(input, "    ")
		require.EqualValues(t, expected, got)
	})

	t.Run("empty line in between", func(t *testing.T) {
		input := `abc

def
`
		expected := `    abc
    
    def
`
		got := AddLinePrefix(input, "    ")
		require.EqualValues(t, expected, got)
	})
}

func Test_AddIndent(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := AddIndent("", "")
		require.EqualValues(t, "", got)
	})

	t.Run("empty content", func(t *testing.T) {
		got := AddIndent("", "prefix")
		require.EqualValues(t, "", got)
	})

	t.Run("empty prefix", func(t *testing.T) {
		got := AddIndent("hello", "")
		require.EqualValues(t, "hello", got)
	})

	t.Run("empty prefix with tailing newline for content", func(t *testing.T) {
		got := AddIndent("hello\n", "")
		require.EqualValues(t, "hello\n", got)
	})

	t.Run("no newline at end", func(t *testing.T) {
		input := `abc
def`
		expected := `    abc
    def`

		got := AddIndent(input, "    ")
		require.EqualValues(t, expected, got)
	})

	t.Run("newline at end", func(t *testing.T) {
		input := `abc
def
`
		expected := `    abc
    def
`
		got := AddIndent(input, "    ")
		require.EqualValues(t, expected, got)
	})

	t.Run("empty line in between", func(t *testing.T) {
		input := `abc

def
`
		expected := `    abc

    def
`
		got := AddIndent(input, "    ")
		require.EqualValues(t, expected, got)
	})
}

func TestStrings_NewStringsService(t *testing.T) {
	service := NewStringsService()
	require.NotNil(t, service)
	require.IsType(t, &StringsService{}, service)
}

func TestStrings_Strings(t *testing.T) {
	service := Strings()
	require.NotNil(t, service)
	require.IsType(t, &StringsService{}, service)
}

func TestStrings_FirstCharToUpper(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"a", "A"},
		{"abc", "Abc"},
		{"ABC", "ABC"},
		{"123", "123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := FirstCharToUpper(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_CountLines(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"test", 1},
		{"test\n", 2},
		{"test\ntest2", 2},
		{"test\ntest2\n", 3},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := CountLines(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_EnsureEndsWithLineBreak(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "\n"},
		{"test", "test\n"},
		{"test\n", "test\n"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := EnsureEndsWithLineBreak(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_EnsureEndsWithExactlyOneLineBreak(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "\n"},
		{"test", "test\n"},
		{"test\n", "test\n"},
		{"test\n\n", "test\n"},
		{"test\n\n\n", "test\n"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := EnsureEndsWithExactlyOneLineBreak(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_RemoveTailingNewline(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"test", "test"},
		{"test\n", "test"},
		{"test\n\n", "test\n"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := RemoveTailingNewline(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_GetFirstLineWithoutCommentAndTrimSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"test", "test"},
		{"# comment\ntest", "test"},
		{"// comment\n\ttest", "test"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := GetFirstLineWithoutCommentAndTrimSpace(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_GetNumberOfLinesWithPrefix(t *testing.T) {
	tests := []struct {
		input    string
		prefix   string
		trim     bool
		expected int
	}{
		{"", "", false, 0},
		{"test", "test", false, 1},
		{"test2\ntest", "test", false, 2},
		{"test\ntest", "test", false, 2},
		{"test\n test", "test", true, 2},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			result := GetNumberOfLinesWithPrefix(tt.input, tt.prefix, tt.trim)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_HasAtLeastOnePrefix(t *testing.T) {
	tests := []struct {
		toCheck  string
		prefixes []string
		expected bool
	}{
		{"", []string{}, false},
		{"test", []string{"pre"}, false},
		{"pretest", []string{"pre"}, true},
		{"test", []string{"pre", "te"}, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			result := HasAtLeastOnePrefix(tt.toCheck, tt.prefixes)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_HasPrefixAndSuffix(t *testing.T) {
	tests := []struct {
		input    string
		prefix   string
		suffix   string
		expected bool
	}{
		{"test", "pre", "suf", false},
		{"pretestsuf", "pre", "suf", true},
		{"pretest", "pre", "suf", false},
		{"testsuf", "pre", "suf", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			result := HasPrefixAndSuffix(tt.input, tt.prefix, tt.suffix)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_SplitAndGetLastElement(t *testing.T) {
	tests := []struct {
		input    string
		token    string
		expected string
	}{
		{"", "/", ""},
		{"test", "/", "test"},
		{"a/b/c", "/", "c"},
		{"a/b/c/d", "/", "d"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := SplitAndGetLastElement(tt.input, tt.token)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_SplitAtSpacesAndRemoveEmptyStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", []string{}},
		{"test", []string{"test"}},
		{"a b c", []string{"a", "b", "c"}},
		{"a  b", []string{"a", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := SplitAtSpacesAndRemoveEmptyStrings(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_SplitFirstLineAndContent(t *testing.T) {
	tests := []struct {
		input      string
		expLine    string
		expContent string
	}{
		{"", "", ""},
		{"test", "test", ""},
		{"test\ncontent", "test", "content"},
		{"line1\nline2\nline3", "line1", "line2\nline3"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			line, content := SplitFirstLineAndContent(tt.input)
			require.Equal(t, tt.expLine, line)
			require.Equal(t, tt.expContent, content)
		})
	}
}

func TestStrings_ToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"test", "Test"},
		{"hello world", "HelloWorld"},
		{"go lang", "GoLang"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToPascalCase(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_ToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"test", "test"},
		{"hello world", "hello_world"},
		{"Go Lang", "go_lang"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToSnakeCase(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_MustGetAsKeyValues(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]string
	}{
		{"", map[string]string{}},
		{"key=value", map[string]string{"key": "value"}},
		{"key1=value1\nkey2=value2", map[string]string{"key1": "value1", "key2": "value2"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := MustGetAsKeyValues(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStrings_MustGetAsKeyValues_Panic(t *testing.T) {
	require.Panics(t, func() {
		MustGetAsKeyValues("invalid line without delimiter")
	})
}

func TestStrings_MustGetValueAsString(t *testing.T) {
	input := "key1=value1\nkey2=value2"
	result := MustGetValueAsString(input, "key1")
	require.Equal(t, "value1", result)
}

func TestStrings_MustGetValueAsString_Panic(t *testing.T) {
	require.Panics(t, func() {
		MustGetValueAsString("key=value", "nonexistent")
	})
}

func TestStrings_MustGetValueAsInt(t *testing.T) {
	input := "key1=123\nkey2=456"
	result := MustGetValueAsInt(input, "key1")
	require.Equal(t, 123, result)
}

func TestStrings_MustGetValueAsInt_Panic(t *testing.T) {
	require.Panics(t, func() {
		MustGetValueAsInt("key=notanumber", "key")
	})
}

func TestStrings_MustHexStringToBytes(t *testing.T) {
	result := MustHexStringToBytes("48656c6c6f")
	require.Equal(t, []byte("Hello"), result)
}

func TestStrings_MustHexStringToBytes_Panic(t *testing.T) {
	require.Panics(t, func() {
		MustHexStringToBytes("invalid hex")
	})
}

func TestStrings_MustMatchesRegex(t *testing.T) {
	result := MustMatchesRegex("test123", "test\\d+")
	require.True(t, result)
}

func TestStrings_MustMatchesRegex_Panic(t *testing.T) {
	require.Panics(t, func() {
		MustMatchesRegex("test", "[invalid")
	})
}

func TestStrings_ContainsAllIgnoreCase(t *testing.T) {
	t.Run("all substrings found", func(t *testing.T) {
		result := ContainsAllIgnoreCase("Hello World", []string{"hello", "world"})
		require.True(t, result)
	})

	t.Run("not all substrings found", func(t *testing.T) {
		result := ContainsAllIgnoreCase("Hello World", []string{"hello", "foo"})
		require.False(t, result)
	})

	t.Run("empty substrings", func(t *testing.T) {
		result := ContainsAllIgnoreCase("Hello World", []string{})
		require.True(t, result)
	})

	t.Run("nil substrings", func(t *testing.T) {
		result := ContainsAllIgnoreCase("Hello World", nil)
		require.True(t, result)
	})
}

func TestStrings_ContainsCommentOnly(t *testing.T) {
	t.Run("comment only", func(t *testing.T) {
		result := ContainsCommentOnly("// comment")
		require.True(t, result)
	})

	t.Run("content with comment", func(t *testing.T) {
		result := ContainsCommentOnly("code\n// comment")
		require.False(t, result)
	})

	t.Run("empty string", func(t *testing.T) {
		result := ContainsCommentOnly("")
		require.False(t, result)
	})

	t.Run("whitespace only", func(t *testing.T) {
		result := ContainsCommentOnly("   ")
		require.False(t, result)
	})
}

func TestStrings_RemoveCommentMarkers(t *testing.T) {
	t.Run("with comment markers", func(t *testing.T) {
		result := RemoveCommentMarkers("// comment\n// line2")
		require.Equal(t, "comment\nline2", result)
	})

	t.Run("empty string", func(t *testing.T) {
		result := RemoveCommentMarkers("")
		require.Equal(t, "", result)
	})
}

func TestStrings_RemoveCommentsAndTrimSpace(t *testing.T) {
	t.Run("with comments", func(t *testing.T) {
		result := RemoveCommentsAndTrimSpace("code\n// comment")
		require.Equal(t, "code", result)
	})

	t.Run("empty string", func(t *testing.T) {
		result := RemoveCommentsAndTrimSpace("")
		require.Equal(t, "", result)
	})
}

func TestStrings_RemoveSurroundingQuotationMarks(t *testing.T) {
	t.Run("with quotes", func(t *testing.T) {
		result := RemoveSurroundingQuotationMarks("\"hello\"")
		require.Equal(t, "hello", result)
	})

	t.Run("without quotes", func(t *testing.T) {
		result := RemoveSurroundingQuotationMarks("hello")
		require.Equal(t, "hello", result)
	})

	t.Run("empty string", func(t *testing.T) {
		result := RemoveSurroundingQuotationMarks("")
		require.Equal(t, "", result)
	})

	t.Run("only opening quote", func(t *testing.T) {
		result := RemoveSurroundingQuotationMarks("\"hello")
		require.Equal(t, "\"hello", result)
	})
}

func TestStrings_RepeatReplaceAll(t *testing.T) {
	t.Run("multiple replacements", func(t *testing.T) {
		result := RepeatReplaceAll("aaa", "a", "b")
		require.Equal(t, "bbb", result)
	})

	t.Run("no match", func(t *testing.T) {
		result := RepeatReplaceAll("hello", "x", "y")
		require.Equal(t, "hello", result)
	})

	t.Run("empty input", func(t *testing.T) {
		result := RepeatReplaceAll("", "a", "b")
		require.Equal(t, "", result)
	})

	t.Run("empty search", func(t *testing.T) {
		result := RepeatReplaceAll("hello", "", "x")
		require.Equal(t, "hello", result)
	})
}

func TestStrings_TrimAllLeadingNewLines(t *testing.T) {
	result := TrimAllLeadingNewLines("\n\n\nhello")
	require.Equal(t, "hello", result)
}

func TestStrings_TrimAllPrefix(t *testing.T) {
	t.Run("multiple prefixes", func(t *testing.T) {
		result := TrimAllPrefix("xxxhello", "x")
		require.Equal(t, "hello", result)
	})

	t.Run("no prefix", func(t *testing.T) {
		result := TrimAllPrefix("hello", "x")
		require.Equal(t, "hello", result)
	})

	t.Run("empty input", func(t *testing.T) {
		result := TrimAllPrefix("", "x")
		require.Equal(t, "", result)
	})

	t.Run("empty prefix", func(t *testing.T) {
		result := TrimAllPrefix("hello", "")
		require.Equal(t, "hello", result)
	})
}

func TestStrings_TrimAllSuffix(t *testing.T) {
	t.Run("multiple suffixes", func(t *testing.T) {
		result := TrimAllSuffix("helloxxx", "x")
		require.Equal(t, "hello", result)
	})

	t.Run("no suffix", func(t *testing.T) {
		result := TrimAllSuffix("hello", "x")
		require.Equal(t, "hello", result)
	})

	t.Run("empty input", func(t *testing.T) {
		result := TrimAllSuffix("", "x")
		require.Equal(t, "", result)
	})

	t.Run("empty suffix", func(t *testing.T) {
		result := TrimAllSuffix("hello", "")
		require.Equal(t, "hello", result)
	})
}

func TestStrings_TrimAllTailingNewLines(t *testing.T) {
	result := TrimAllTailingNewLines("hello\n\n\n")
	require.Equal(t, "hello", result)
}

func TestStrings_TrimPrefixAndSuffix(t *testing.T) {
	result := TrimPrefixAndSuffix("<hello>", "<", ">")
	require.Equal(t, "hello", result)
}

func TestStrings_TrimSpaceForEveryLine(t *testing.T) {
	t.Run("reconstructs lines", func(t *testing.T) {
		result := TrimSpaceForEveryLine("hello\nworld")
		require.Equal(t, "hello\nworld", result)
	})
}

func TestStrings_TrimSpacesRight(t *testing.T) {
	result := TrimSpacesRight("hello  \t\n")
	require.Equal(t, "hello", result)
}

func TestStrings_TrimSuffixAndSpace(t *testing.T) {
	t.Run("suffix at end", func(t *testing.T) {
		result := TrimSuffixAndSpace("hello!", "!")
		require.Equal(t, "hello", result)
	})

	t.Run("suffix with trailing space", func(t *testing.T) {
		result := TrimSuffixAndSpace("hello!  ", "!")
		require.Equal(t, "hello!", result)
	})
}

func TestStrings_TrimSuffixUntilAbsent(t *testing.T) {
	t.Run("multiple suffixes", func(t *testing.T) {
		result := TrimSuffixUntilAbsent("hello!!!", "!")
		require.Equal(t, "hello", result)
	})

	t.Run("no suffix", func(t *testing.T) {
		result := TrimSuffixUntilAbsent("hello", "!")
		require.Equal(t, "hello", result)
	})
}

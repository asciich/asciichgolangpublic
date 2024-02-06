package asciichgolangpublic

import (
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				firstLine := Strings().GetFirstLine(tt.input)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				firstLine := Strings().GetFirstLineAndTrimSpace(tt.input)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				ensuredLineBreak := Strings().EnsureEndsWithExactlyOneLineBreak(tt.input)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				ensuredLineBreak := Strings().RemoveTailingNewline(tt.input)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				firstCharUppercased := Strings().EnsureFirstCharUppercase(tt.input)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				firstCharUppercased := Strings().EnsureFirstCharLowercase(tt.input)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				commentsRemoved := Strings().RemoveComments(tt.input)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				filled := Strings().RightFillWithSpaces(tt.input, tt.fillLenght)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedHasPrefix,
					Strings().HasPrefixIgnoreCase(tt.input, tt.prefix),
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedTrimmed,
					Strings().TrimPrefixIgnoreCase(tt.input, tt.prefix),
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedFirstCharLowerCase,
					Strings().IsFirstCharLowerCase(tt.input),
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedFirstCharLowerCase,
					Strings().IsFirstCharUpperCase(tt.input),
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedWords,
					Strings().SplitWords(tt.input),
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedIsComment,
					Strings().IsComment(tt.input),
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedOutput,
					Strings().TrimSpacesLeft(tt.input),
				)
			},
		)
	}
}

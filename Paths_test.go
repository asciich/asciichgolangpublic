package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathsIsRelativePath(t *testing.T) {

	tests := []struct {
		path               string
		expectedIsRelative bool
	}{
		{"", false},
		{"this", true},
		{"this/is/relative", true},
		{"/this/is/absoute", false},
		{"/", false},
		{"c:\\", false},
		{"c:\\Users", false},
		{"C:\\", false},
		{"C:\\Users", false},
		{"d:\\", false},
		{"d:\\Users", false},
		{"D:\\", false},
		{"D:\\Users", false},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedIsRelative,
					Paths().IsRelativePath(tt.path),
				)
			},
		)
	}
}

func TestPathsIsAbsolutePath(t *testing.T) {

	tests := []struct {
		path               string
		expectedIsRelative bool
	}{
		{"", false},
		{"this", false},
		{"this/is/relative", false},
		{"/this/is/absoute", true},
		{"/", true},
		{"c:\\", true},
		{"c:\\Users", true},
		{"C:\\", true},
		{"C:\\Users", true},
		{"d:\\", true},
		{"d:\\Users", true},
		{"D:\\", true},
		{"D:\\Users", true},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedIsRelative,
					Paths().IsAbsolutePath(tt.path),
				)
			},
		)
	}
}

func TestPaths_MatchBaseNamePattern(t *testing.T) {
	tests := []struct {
		pathFilterOptions PathFilterOptions
		expectedFileList  []string
	}{
		{&ListFileOptions{}, []string{"a.html", "a.txt", "b.html", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{}}, []string{"a.html", "a.txt", "b.html", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{".*"}}, []string{"a.html", "a.txt", "b.html", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{"^txt$"}}, []string{}},
		{&ListFileOptions{MatchBasenamePattern: []string{"txt$"}}, []string{"a.txt", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{"txt"}}, []string{"a.txt", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{"txt", "xt"}}, []string{"a.txt", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{".*\\.txt"}}, []string{"a.txt", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{".*\\.txt$"}}, []string{"a.txt", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{"^.*\\.txt$"}}, []string{"a.txt", "b.txt"}},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				input := []string{
					"a.txt",
					"b.txt",
					"a.html",
					"b.html",
				}

				assert.EqualValues(
					tt.expectedFileList,
					Paths().MustFilterPaths(input, tt.pathFilterOptions),
				)
			},
		)
	}
}

func TestPaths_MatchBaseNamePattern_recursive(t *testing.T) {
	tests := []struct {
		pathFilterOptions PathFilterOptions
		expectedFileList  []string
	}{
		{&ListFileOptions{}, []string{"a.html", "a.txt", "abc/a.html", "abc/a.txt", "abc/b.html", "abc/b.txt", "b.html", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{}}, []string{"a.html", "a.txt", "abc/a.html", "abc/a.txt", "abc/b.html", "abc/b.txt", "b.html", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{".*"}}, []string{"a.html", "a.txt", "abc/a.html", "abc/a.txt", "abc/b.html", "abc/b.txt", "b.html", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{"^txt$"}}, []string{}},
		{&ListFileOptions{MatchBasenamePattern: []string{"txt$"}}, []string{"a.txt", "abc/a.txt", "abc/b.txt", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{"txt"}}, []string{"a.txt", "abc/a.txt", "abc/b.txt", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{"txt", "xt"}}, []string{"a.txt", "abc/a.txt", "abc/b.txt", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{".*\\.txt"}}, []string{"a.txt", "abc/a.txt", "abc/b.txt", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{".*\\.txt$"}}, []string{"a.txt", "abc/a.txt", "abc/b.txt", "b.txt"}},
		{&ListFileOptions{MatchBasenamePattern: []string{"^.*\\.txt$"}}, []string{"a.txt", "abc/a.txt", "abc/b.txt", "b.txt"}},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				input := []string{
					"a.txt",
					"b.txt",
					"a.html",
					"b.html",
					"abc/a.txt",
					"abc/b.txt",
					"abc/a.html",
					"abc/b.html",
				}

				assert.EqualValues(
					tt.expectedFileList,
					Paths().MustFilterPaths(input, tt.pathFilterOptions),
				)
			},
		)
	}
}

func TestPaths_ExcludeBasenamePattern_recursive(t *testing.T) {
	tests := []struct {
		pathFilterOptions PathFilterOptions
		expectedFileList  []string
	}{
		{&ListFileOptions{}, []string{"a.html", "a.txt", "abc/a.html", "abc/a.txt", "abc/b.html", "abc/b.txt", "b.html", "b.txt"}},
		{&ListFileOptions{ExcludeBasenamePattern: []string{".*\\.html"}}, []string{"a.txt", "abc/a.txt", "abc/b.txt", "b.txt"}},
		{&ListFileOptions{ExcludeBasenamePattern: []string{".*\\.txt"}}, []string{"a.html", "abc/a.html", "abc/b.html", "b.html"}},
		{&ListFileOptions{ExcludeBasenamePattern: []string{".*\\.txt"}, MatchBasenamePattern: []string{"^a.*"}}, []string{"a.html", "abc/a.html"}},
		{&ListFileOptions{ExcludeBasenamePattern: []string{".*\\.txt", "^b.*"}, MatchBasenamePattern: []string{"^a.*"}}, []string{"a.html", "abc/a.html"}},
		{&ListFileOptions{ExcludeBasenamePattern: []string{".*\\.txt", "^b.*"}}, []string{"a.html", "abc/a.html"}},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				input := []string{
					"a.txt",
					"b.txt",
					"a.html",
					"b.html",
					"abc/a.txt",
					"abc/b.txt",
					"abc/a.html",
					"abc/b.html",
				}

				assert.EqualValues(
					tt.expectedFileList,
					Paths().MustFilterPaths(input, tt.pathFilterOptions),
				)
			},
		)
	}
}

func TestPaths_ExcludeWholepathPattern_recursive(t *testing.T) {
	tests := []struct {
		pathFilterOptions PathFilterOptions
		expectedFileList  []string
	}{
		{&ListFileOptions{}, []string{"a.html", "a.txt", "abc/a.html", "abc/a.txt", "abc/b.html", "abc/b.txt", "b.html", "b.txt"}},
		{&ListFileOptions{ExcludePatternWholepath: []string{"c/a"}}, []string{"a.html", "a.txt", "abc/b.html", "abc/b.txt", "b.html", "b.txt"}},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				input := []string{
					"a.txt",
					"b.txt",
					"a.html",
					"b.html",
					"abc/a.txt",
					"abc/b.txt",
					"abc/a.html",
					"abc/b.html",
				}

				assert.EqualValues(
					tt.expectedFileList,
					Paths().MustFilterPaths(input, tt.pathFilterOptions),
				)
			},
		)
	}
}

func TestPaths_GetRelativePathTo(t *testing.T) {
	tests := []struct {
		input          string
		relativeTo     string
		expectedOutput string
	}{
		{"/bin/bash", "/", "bin/bash"},
		{"/bin/bash", "/bin", "bash"},
		{"/bin/bash", "/bin/", "bash"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedOutput,
					Paths().MustGetRelativePathTo(tt.input, tt.relativeTo),
				)
			},
		)
	}
}


func TestPaths_GetRelativePathsTo(t *testing.T) {
	tests := []struct {
		input          []string
		relativeTo     string
		expectedOutput []string
	}{
		{[]string{"/bin/bash", "/bin/cat"}, "/", []string{"bin/bash", "bin/cat"}},
		{[]string{"/bin/bash", "/bin/cat"}, "/bin", []string{"bash", "cat"}},
		{[]string{"/bin/bash", "/bin/cat"}, "/bin/", []string{"bash", "cat"}},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedOutput,
					Paths().MustGetRelativePathsTo(tt.input, tt.relativeTo),
				)
			},
		)
	}
}

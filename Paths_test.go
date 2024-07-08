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

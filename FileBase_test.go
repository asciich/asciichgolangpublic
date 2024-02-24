package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileBase(t *testing.T) {
	tests := []struct {
		command        []string
		expectedOutput string
	}{
		{[]string{"echo", "hello"}, "hello\n"},
		{[]string{"echo", "hello world"}, "hello world\n"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				fileBase := FileBase{}

				parent, err := fileBase.GetParentFileForBaseClass()
				assert.Nil(parent)
				assert.ErrorIs(err, ErrFileBaseParentNotSet)
				assert.ErrorIs(err, ErrTracedError)
			},
		)
	}
}

package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectoryBaseSetAndGetParentDirectory(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				directoryBase := NewDirectoryBase()
				directory := TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose)
				defer directory.Delete(verbose)

				directoryBase.MustSetParentDirectoryForBaseClass(directory)

				assert.EqualValues(
					directoryBase.MustGetParentDirectoryForBaseClass(),
					directory,
				)
			},
		)
	}
}

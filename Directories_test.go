package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectoriesCreateLocalDirectoryByPath(t *testing.T) {

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

				var directory Directory = TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose)
				defer directory.Delete(verbose)

				assert.True(directory.MustExists(verbose))

				for i := 0; i < 2; i++ {
					directory.MustDelete(verbose)
					assert.False(directory.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					createdDir := Directories().MustCreateLocalDirectoryByPath(directory.MustGetLocalPath(), verbose)
					assert.True(directory.MustExists(verbose))
					assert.True(createdDir.MustExists(verbose))
				}

			},
		)
	}
}

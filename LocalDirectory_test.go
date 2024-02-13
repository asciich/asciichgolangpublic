package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalDirectoryExists(t *testing.T) {

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

				assert.True(directory.MustExists())

				for i := 0; i < 2; i++ {
					directory.MustDelete(verbose)
					assert.False(directory.MustExists())
				}

				for i := 0; i < 2; i++ {
					directory.MustCreate(verbose)
					assert.True(directory.MustExists())
				}

				for i := 0; i < 2; i++ {
					directory.MustDelete(verbose)
					assert.False(directory.MustExists())
				}

			},
		)
	}
}

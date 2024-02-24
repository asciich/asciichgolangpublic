package asciichgolangpublic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserGetHomeDirectory(t *testing.T) {

	tests := []struct {
		stringName string
	}{
		{"varName"},
		{"AnoterVarName"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.True(
					strings.HasPrefix(
						Users().MustGetHomeDirectoryAsString(),
						"/home/",
					),
				)
			},
		)
	}
}

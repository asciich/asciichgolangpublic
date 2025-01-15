package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestVersionDateVersionGetAsString(t *testing.T) {
	tests := []struct {
		testmessage string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				var version Version = Versions().MustGetNewDateVersion()
				versionString := version.MustGetAsString()

				assert.True(Versions().IsVersionString(versionString))
			},
		)
	}
}

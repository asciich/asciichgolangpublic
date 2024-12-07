package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64_encodeAndDecode(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"hello world"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				encoded := Base64().MustEncodeStringAsString(tt.input)

				assert.EqualValues(
					tt.input,
					Base64().MustDecodeStringAsString(encoded),
				)
			},
		)
	}
}

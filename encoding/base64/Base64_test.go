package base64

import (
	"fmt"
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
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				encoded := MustEncodeStringAsString(tt.input)

				assert.EqualValues(
					tt.input,
					MustDecodeStringAsString(encoded),
				)
			},
		)
	}
}

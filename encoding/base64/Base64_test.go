package base64

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
				require := require.New(t)

				encoded := MustEncodeStringAsString(tt.input)

				require.EqualValues(
					tt.input,
					MustDecodeStringAsString(encoded),
				)
			},
		)
	}
}

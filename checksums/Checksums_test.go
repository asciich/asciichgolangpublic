package checksums

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChecksumsGetSha256SumFromString(t *testing.T) {
	tests := []struct {
		input            string
		expectedChecksum string
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"hello world", "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := GetSha256SumFromString(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

func TestChecksumsGetSha256SumFromBytes(t *testing.T) {
	tests := []struct {
		input            []byte
		expectedChecksum string
	}{
		{[]byte(""), "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{[]byte("hello world"), "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := GetSha256SumFromBytes(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

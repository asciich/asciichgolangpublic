package bytesutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesParseSizeStringAsInt64(t *testing.T) {
	tests := []struct {
		stringToParse     string
		expectedSizeBytes int64
	}{
		{"1", 1},
		{"10", 10},
		{"100", 100},
		{"1kB", 1000},
		{"1KB", 1024},
		{"1KiB", 1024},
		{"0.5KiB", 1024 / 2},
		{"2KiB", 2048},
		{"1MB", 1024 * 1024},
		{"1MiB", 1024 * 1024},
		{"1GB", 1024 * 1024 * 1024},
		{"1GiB", 1024 * 1024 * 1024},
		{"1TB", 1024 * 1024 * 1024 * 1024},
		{"1TiB", 1024 * 1024 * 1024 * 1024},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				sizeBytes := MustParseSizeStringAsInt64(tt.stringToParse)
				require.EqualValues(tt.expectedSizeBytes, sizeBytes)
			},
		)
	}
}

func TestGetSizeAsHumanReadableString(t *testing.T) {
	tests := []struct {
		sizeToConvert      int64
		expectedSizeString string
	}{
		{1, "1"},
		{10, "10"},
		{100, "100"},
		{1000, "1000"},
		{1024, "1KiB"},
		{1024 * 1024, "1MiB"},
		{1024 * 1024 * 1024, "1GiB"},
		{1024 * 1024 * 1024 * 1024, "1TiB"},
		{512 * 1953525168, "931.51GiB"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				sizeString := MustGetSizeAsHumanReadableString(tt.sizeToConvert)
				require.EqualValues(tt.expectedSizeString, sizeString)
			},
		)
	}
}

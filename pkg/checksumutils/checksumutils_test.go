package checksumutils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/checksumutils"
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
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetSha256SumFromString(tt.input)
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
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetSha256SumFromBytes(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

func TestChecksumsGetSha1SumFromString(t *testing.T) {
	tests := []struct {
		input            string
		expectedChecksum string
	}{
		{"", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{"hello world", "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetSha1SumFromString(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

func TestChecksumsGetSha1SumFromBytes(t *testing.T) {
	tests := []struct {
		input            []byte
		expectedChecksum string
	}{
		{[]byte(""), "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{[]byte("hello world"), "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetSha1SumFromBytes(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

func TestChecksumsGetMD5SumFromString(t *testing.T) {
	tests := []struct {
		input            string
		expectedChecksum string
	}{
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"hello world", "5eb63bbbe01eeed093cb22bb8f5acdc3"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetMD5SumFromString(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

func TestChecksumsGetMD5SumFromBytes(t *testing.T) {
	tests := []struct {
		input            []byte
		expectedChecksum string
	}{
		{[]byte(""), "d41d8cd98f00b204e9800998ecf8427e"},
		{[]byte("hello world"), "5eb63bbbe01eeed093cb22bb8f5acdc3"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetMD5SumFromBytes(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

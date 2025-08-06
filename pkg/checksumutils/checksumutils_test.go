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

func TestChecksumsGetSha512SumFromString(t *testing.T) {
	tests := []struct {
		input            string
		expectedChecksum string
	}{
		{"", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
		{"hello world", "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetSha512SumFromString(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

func TestChecksumsGetSha512SumFromBytes(t *testing.T) {
	tests := []struct {
		input            []byte
		expectedChecksum string
	}{
		{[]byte(""), "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
		{[]byte("hello world"), "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetSha512SumFromBytes(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

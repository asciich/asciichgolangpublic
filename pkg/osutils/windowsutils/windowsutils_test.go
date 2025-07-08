package windowsutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/windowsutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestWindowsIsRunningOnWindows(t *testing.T) {
	tests := []struct {
		testmessage string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				require.False(windowsutils.IsRunningOnWindows())
			},
		)
	}
}

func TestWindowsDecodeAsString(t *testing.T) {
	tests := []struct {
		inputUtf16     []byte
		expectedOutput string
	}{
		{nil, ""},
		{[]byte{}, ""},
		{[]byte{0x48, 0x00}, "H"},
		{[]byte{0x48, 0x00, 0x65, 0x00}, "He"},
		{[]byte{0x48, 0x00, 0x65, 0x00, 0x6c, 0x00, 0x6c, 0x00, 0x6f, 0x00}, "Hello"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				out, err := windowsutils.DecodeAsString(tt.inputUtf16)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedOutput, out)
			},
		)
	}
}

package utf16utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/encodingutils/utf16utils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestUTF16DecodeAsString(t *testing.T) {
	tests := []struct {
		inputUtf16     []byte
		expectedOutput string
	}{
		{nil, ""},
		{[]byte{}, ""},
		{[]byte{0x48, 0x00}, "H"},
		{[]byte{0x48, 0x00, 0x65, 0x00}, "He"},
		{[]byte{0x48, 0x00, 0x65, 0x00, 0x6c, 0x00, 0x6c, 0x00, 0x6f, 0x00}, "Hello"},

		// Do not convert already UTF8 strings:
		{[]byte("h"), "h"},
		{[]byte("hello"), "hello"},
		{[]byte(" hello"), " hello"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				decoded, err := utf16utils.DecodeAsString(tt.inputUtf16)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedOutput, decoded)
			},
		)
	}
}

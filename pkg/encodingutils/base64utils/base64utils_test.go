package base64utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/encodingutils/base64utils"
)

func TestBase64_encodeAndDecode(t *testing.T) {
	t.Run("hello world", func(t *testing.T) {
		encoded, err := base64utils.EncodeStringAsString("hello world")
		require.NoError(t, err)

		decoded, err := base64utils.DecodeStringAsString(encoded)
		require.NoError(t, err)

		require.EqualValues(t, "hello world", decoded)
	},
	)
}

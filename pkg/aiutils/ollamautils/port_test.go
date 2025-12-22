package ollamautils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/ollamautils"
)

func Test_GetDefaultPort(t *testing.T) {
	require.EqualValues(t, 11434, ollamautils.GetDefaultPort())
}
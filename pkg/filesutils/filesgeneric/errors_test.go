package filesgeneric_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesgeneric"
)

func Test_IsErrFileNotFound(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.False(t, filesgeneric.IsErrFileNotFound(nil))
	})

	t.Run("another error", func(t *testing.T) {
		require.False(t, filesgeneric.IsErrFileNotFound(fmt.Errorf("another error")))
	})

	t.Run("direct", func(t *testing.T) {
		require.True(t, filesgeneric.IsErrFileNotFound(filesgeneric.ErrFileNotFound))
	})

	t.Run("wrapped", func(t *testing.T) {
		require.True(t, filesgeneric.IsErrFileNotFound(fmt.Errorf("This is wrapping: %w", filesgeneric.ErrFileNotFound)))
	})

	t.Run("The default os.ErrNotExist returned by os.Open is ErrFileNotFound as well", func(t *testing.T) {
		_, err := os.Open("/this/file/does/not/exist")
		require.True(t, filesgeneric.IsErrFileNotFound(err))
	})
}

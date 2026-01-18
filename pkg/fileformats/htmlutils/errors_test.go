package htmlutils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/htmlutils"
)

func Test_IsErrNoHtmlBodyFound(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.False(t, htmlutils.IsErrNoHtmlBodyFound(nil))
	})

	t.Run("another error", func(t *testing.T) {
		require.False(t, htmlutils.IsErrNoHtmlBodyFound(fmt.Errorf("this is another error")))
	})

	t.Run("error wrapped", func(t *testing.T) {
		require.True(t, htmlutils.IsErrNoHtmlBodyFound(fmt.Errorf("%w", htmlutils.ErrNoHtmlBodyFound)))
	})
}

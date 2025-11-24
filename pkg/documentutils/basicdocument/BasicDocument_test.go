package basicdocument_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/basicdocument"
)

func Test_AddTitleByString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		require.Error(t, basicdocument.NewBasicDocument().AddTitleByString(""))
	})

	t.Run("valid title", func(t *testing.T) {
		require.NoError(t, basicdocument.NewBasicDocument().AddTitleByString("valid title"))
	})
}

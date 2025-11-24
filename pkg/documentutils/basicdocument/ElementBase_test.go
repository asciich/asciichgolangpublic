package basicdocument

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SetAndGetPlainText(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		require.Error(t, NewElementBase().SetPlainText(""))
	})

	t.Run("valid text", func(t *testing.T) {
		e := NewElementBase()
		require.NoError(t, e.SetPlainText("valid title"))

		p, err := e.GetPlainText()
		require.NoError(t, err)
		require.EqualValues(t, "valid title", p)
	})
}

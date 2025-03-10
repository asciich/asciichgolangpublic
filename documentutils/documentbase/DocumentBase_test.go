package documentbase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AddTitleByString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		require.Error(t, NewDocumentBase().AddTitleByString(""))
	})

	t.Run("valid title", func(t *testing.T) {
		d := NewDocumentBase()
		require.NoError(t, d.AddTitleByString("valid title"))

		// firstElement := d.GetElements()[0]
		
	})
}

func Test_RenderAsString(t *testing.T) {
	// The base class itself does not implemnet any renderer.
	// Therefore it will always return an error.

	rendered, err := NewDocumentBase().RenderAsString()
	require.Error(t, err)
	require.EqualValues(t, "", rendered)
}
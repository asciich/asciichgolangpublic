package markdowndocument

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RenderEmpty(t *testing.T) {
	rendered, err := NewMarkDownDocument().RenderAsString()
	require.NoError(t, err)
	require.EqualValues(t, "\n", rendered)
}

func Test_Render(t *testing.T) {
	t.Run("Only title", func(t *testing.T) {
		d := NewMarkDownDocument()
		require.NoError(t, d.AddTitleByString("example title"))

		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n", rendered)
	})

	t.Run("Only subtitle", func(t *testing.T) {
		d := NewMarkDownDocument()
		require.NoError(t, d.AddSubTitleByString("example title"))

		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "## example title\n", rendered)
	})

	t.Run("Only subsubtitle", func(t *testing.T) {
		d := NewMarkDownDocument()
		require.NoError(t, d.AddSubSubTitleByString("example title"))

		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "### example title\n", rendered)
	})

	t.Run("Only subsubsubtitle", func(t *testing.T) {
		d := NewMarkDownDocument()
		require.NoError(t, d.AddSubSubSubTitleByString("example title"))

		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "#### example title\n", rendered)
	})


	t.Run("title and text", func(t *testing.T) {
		d := NewMarkDownDocument()
		require.NoError(t, d.AddTitleByString("example title"))
		require.NoError(t, d.AddTextByString("example text."))

		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\nexample text.\n", rendered)
	})

	t.Run("title and two text", func(t *testing.T) {
		d := NewMarkDownDocument()
		require.NoError(t, d.AddTitleByString("example title"))
		require.NoError(t, d.AddTextByString("example text."))
		require.NoError(t, d.AddTextByString("example text2."))

		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\nexample text.\n\nexample text2.\n", rendered)
	})

	t.Run("title, text, title, text", func(t *testing.T) {
		d := NewMarkDownDocument()
		require.NoError(t, d.AddTitleByString("example title"))
		require.NoError(t, d.AddTextByString("example text."))
		require.NoError(t, d.AddTitleByString("example title2"))
		require.NoError(t, d.AddTextByString("example text2."))

		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\nexample text.\n\n# example title2\n\nexample text2.\n", rendered)
	})

}

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

	t.Run("title, table", func(t *testing.T) {
		d := NewMarkDownDocument()
		require.NoError(t, d.AddTitleByString("example title"))
		spreadsheet, err := d.AddTable()
		require.NoError(t, err)

		require.NoError(t, spreadsheet.SetColumnTitles([]string{"col1", "col nr 2"}))

		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\n| col1 | col nr 2 |\n| ---- | -------- |\n", rendered)
	})
	t.Run("title, verbatim, text, verbatim", func(t *testing.T) {
		d := NewMarkDownDocument()
		require.NoError(t, d.AddTitleByString("example title"))
		require.NoError(t, d.AddVerbatimByString("code block 1"))
		require.NoError(t, d.AddTextByString("example text."))
		require.NoError(t, d.AddVerbatimByString("code block 2"))

		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\n```\ncode block 1\n```\n\nexample text.\n\n```\ncode block 2\n```\n", rendered)
	})
}

func Test_ParseFromString(t *testing.T) {
	t.Run("Empty string", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "\n", rendered)
	})

	t.Run("Only title", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("# example title")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n", rendered)
	})

	t.Run("Only subtitle", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("## example subtitle")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "## example subtitle\n", rendered)
	})

	t.Run("Only subsubtitle", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("### example subsubtitle")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "### example subsubtitle\n", rendered)
	})

	t.Run("Only subsubsubtitle", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("#### example subsubsubtitle")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "#### example subsubsubtitle\n", rendered)
	})

	t.Run("Title and text", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("# example title\n\nexample text.")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\nexample text.\n", rendered)
	})

	t.Run("Title and two text paragraphs", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("# example title\n\nexample text.\n\nexample text2.")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\nexample text.\n\nexample text2.\n", rendered)
	})

	t.Run("Multiple headings and text", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("# example title\n\nexample text.\n\n# example title2\n\nexample text2.")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\nexample text.\n\n# example title2\n\nexample text2.\n", rendered)
	})

	t.Run("All heading levels", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("# Title\n## Subtitle\n### SubSubTitle\n#### SubSubSubTitle")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# Title\n\n## Subtitle\n\n### SubSubTitle\n\n#### SubSubSubTitle\n", rendered)
	})

	t.Run("Ignores empty lines", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("\n\n# example title\n\n\nexample text.\n\n")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\nexample text.\n", rendered)
	})

	t.Run("Ignores table lines", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("# Title\n| col1 | col2 |\n| ---- | ---- |\nsome text")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# Title\n\nsome text\n", rendered)
	})

	t.Run("Verbatim block", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("```\ncode block\n```")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "```\ncode block\n```\n", rendered)
	})

	t.Run("Title with verbatim", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("# example title\n\n```\ncode block 1\n```")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\n```\ncode block 1\n```\n", rendered)
	})

	t.Run("Multiple verbatim blocks with text", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("# example title\n\n```\ncode block 1\n```\n\nexample text.\n\n```\ncode block 2\n```")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\n```\ncode block 1\n```\n\nexample text.\n\n```\ncode block 2\n```\n", rendered)
	})

	t.Run("Verbatim with multiple lines", func(t *testing.T) {
		d := NewMarkDownDocument()
		err := d.ParseFromString("```\nline 1\nline 2\nline 3\n```")
		require.NoError(t, err)
		
		rendered, err := d.RenderAsString()
		require.NoError(t, err)
		require.EqualValues(t, "```\nline 1\nline 2\nline 3\n```\n", rendered)
	})
}
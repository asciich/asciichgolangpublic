package markdowndocument_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/basicdocument"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/markdowndocument"
)

func Test_RenderEmpty(t *testing.T) {
	document := basicdocument.NewBasicDocument()

	rendered, err := markdowndocument.RenderAsString(document)
	require.NoError(t, err)
	require.EqualValues(t, "\n", rendered)
}

func Test_Render(t *testing.T) {
	t.Run("Only title", func(t *testing.T) {
		document := basicdocument.NewBasicDocument()
		require.NoError(t, document.AddTitleByString("example title"))

		rendered, err := markdowndocument.RenderAsString(document)
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n", rendered)
	})

	t.Run("Only subtitle", func(t *testing.T) {
		document := basicdocument.NewBasicDocument()
		require.NoError(t, document.AddSubTitleByString("example title"))

		rendered, err := markdowndocument.RenderAsString(document)
		require.NoError(t, err)
		require.EqualValues(t, "## example title\n", rendered)
	})

	t.Run("Only subsubtitle", func(t *testing.T) {
		document := basicdocument.NewBasicDocument()
		require.NoError(t, document.AddSubSubTitleByString("example title"))

		rendered, err := markdowndocument.RenderAsString(document)
		require.NoError(t, err)
		require.EqualValues(t, "### example title\n", rendered)
	})

	t.Run("Only subsubsubtitle", func(t *testing.T) {
		document := basicdocument.NewBasicDocument()
		require.NoError(t, document.AddSubSubSubTitleByString("example title"))

		rendered, err := markdowndocument.RenderAsString(document)
		require.NoError(t, err)
		require.EqualValues(t, "#### example title\n", rendered)
	})

	t.Run("title and text", func(t *testing.T) {
		document := basicdocument.NewBasicDocument()
		require.NoError(t, document.AddTitleByString("example title"))
		require.NoError(t, document.AddTextByString("example text."))

		rendered, err := markdowndocument.RenderAsString(document)
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\nexample text.\n", rendered)
	})

	t.Run("title and two text", func(t *testing.T) {
		document := basicdocument.NewBasicDocument()
		require.NoError(t, document.AddTitleByString("example title"))
		require.NoError(t, document.AddTextByString("example text."))
		require.NoError(t, document.AddTextByString("example text2."))

		rendered, err := markdowndocument.RenderAsString(document)
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\nexample text.\n\nexample text2.\n", rendered)
	})

	t.Run("title, text, title, text", func(t *testing.T) {
		document := basicdocument.NewBasicDocument()
		require.NoError(t, document.AddTitleByString("example title"))
		require.NoError(t, document.AddTextByString("example text."))
		require.NoError(t, document.AddTitleByString("example title2"))
		require.NoError(t, document.AddTextByString("example text2."))
		rendered, err := markdowndocument.RenderAsString(document)
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\nexample text.\n\n# example title2\n\nexample text2.\n", rendered)
	})

	t.Run("title, table", func(t *testing.T) {
		document := basicdocument.NewBasicDocument()
		require.NoError(t, document.AddTitleByString("example title"))
		spreadsheet, err := document.AddTable()
		require.NoError(t, err)

		require.NoError(t, spreadsheet.SetColumnTitles([]string{"col1", "col nr 2"}))

		rendered, err := markdowndocument.RenderAsString(document)
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\n| col1 | col nr 2 |\n| ---- | -------- |\n", rendered)
	})
	t.Run("title, verbatim, text, verbatim", func(t *testing.T) {
		document := basicdocument.NewBasicDocument()
		require.NoError(t, document.AddTitleByString("example title"))
		require.NoError(t, document.AddVerbatimByString("code block 1"))
		require.NoError(t, document.AddTextByString("example text."))
		require.NoError(t, document.AddVerbatimByString("code block 2"))
		rendered, err := markdowndocument.RenderAsString(document)
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\n```\ncode block 1\n```\n\nexample text.\n\n```\ncode block 2\n```\n", rendered)
	})

}

func Test_RenderCodeBlock(t *testing.T) {
	t.Run("title, text, codeblock without langauge", func(t *testing.T) {
		document := basicdocument.NewBasicDocument()
		require.NoError(t, document.AddTitleByString("example title"))
		require.NoError(t, document.AddTextByString("example text."))
		require.NoError(t, document.AddCodeBlockByString("example command", ""))
		rendered, err := markdowndocument.RenderAsString(document)
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\nexample text.\n\n```\nexample command\n```\n", rendered)
	})

	t.Run("title, text, codeblock shell", func(t *testing.T) {
		document := basicdocument.NewBasicDocument()
		require.NoError(t, document.AddTitleByString("example title"))
		require.NoError(t, document.AddTextByString("example text."))
		require.NoError(t, document.AddCodeBlockByString("example command", "shell"))
		rendered, err := markdowndocument.RenderAsString(document)
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\nexample text.\n\n```shell\nexample command\n```\n", rendered)
	})

	t.Run("title, text, codeblock bash", func(t *testing.T) {
		document := basicdocument.NewBasicDocument()
		require.NoError(t, document.AddTitleByString("example title"))
		require.NoError(t, document.AddTextByString("example text."))
		require.NoError(t, document.AddCodeBlockByString("example command", "bash"))
		rendered, err := markdowndocument.RenderAsString(document)
		require.NoError(t, err)
		require.EqualValues(t, "# example title\n\nexample text.\n\n```bash\nexample command\n```\n", rendered)
	})
}

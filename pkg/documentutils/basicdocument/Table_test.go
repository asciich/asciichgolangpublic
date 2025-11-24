package  basicdocument

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func Test_GetSpreadsheetFromTable(t *testing.T) {
	table := mustutils.Must(GetNewTable())
	spreadsheet := mustutils.Must(table.GetSpreadSheet())

	require.True(t, spreadsheet.MustIsEmpty())
}

func Test_RenderMarkDown(t *testing.T) {
	t.Run("only one title", func(t *testing.T) {
		table := mustutils.Must(GetNewTable())
		spreadsheet := mustutils.Must(table.GetSpreadSheet())
		spreadsheet.MustSetColumnTitles([]string{"one title"})

		rendered := table.MustRenderAsMarkdownString()
		require.EqualValues(t, "| one title |\n| --------- |\n", rendered)
	})

	t.Run("one title one entry", func(t *testing.T) {
		table := mustutils.Must(GetNewTable())
		spreadsheet := mustutils.Must(table.GetSpreadSheet())
		spreadsheet.MustSetColumnTitles([]string{"one title"})
		spreadsheet.MustAddRow([]string{"entry"})

		rendered := table.MustRenderAsMarkdownString()
		require.EqualValues(t, "| one title |\n| --------- |\n| entry |\n", rendered)
	})
}

package basicdocument

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
		err := spreadsheet.SetColumnTitles([]string{"one title"})
		require.NoError(t, err)

		rendered := table.MustRenderAsMarkdownString()
		require.EqualValues(t, "| one title |\n| --------- |\n", rendered)
	})

	t.Run("one title one entry", func(t *testing.T) {
		table := mustutils.Must(GetNewTable())
		spreadsheet := mustutils.Must(table.GetSpreadSheet())

		err := spreadsheet.SetColumnTitles([]string{"one title"})
		require.NoError(t, err)

		err = spreadsheet.AddRow([]string{"entry"})
		require.NoError(t, err)

		rendered := table.MustRenderAsMarkdownString()
		require.EqualValues(t, "| one title |\n| --------- |\n| entry |\n", rendered)
	})
}

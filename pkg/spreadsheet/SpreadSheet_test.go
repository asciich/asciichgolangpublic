package spreadsheet

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestSpreadSheetNoRowsAndColumns(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				spreadSheet := NewSpreadSheet()

				numberOfColumns, err := spreadSheet.GetNumberOfColumns()
				require.NoError(t, err)
				require.EqualValues(t, 0, numberOfColumns)

				numberOfRows, err := spreadSheet.GetNumberOfRows()
				require.NoError(t, err)
				require.EqualValues(t, 0, numberOfRows)

				isEmpty, err := spreadSheet.IsEmpty()
				require.NoError(t, err)
				require.True(t, isEmpty)
			},
		)
	}
}

func TestSpreadSheetSetColumnTitles(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				spreadSheet := NewSpreadSheet()

				err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
				require.NoError(t, err)

				numberOfColumns, err := spreadSheet.GetNumberOfColumns()
				require.NoError(t, err)
				require.EqualValues(t, 2, numberOfColumns)

				numberOfRows, err := spreadSheet.GetNumberOfRows()
				require.NoError(t, err)
				require.EqualValues(t, 0, numberOfRows)
			},
		)
	}
}

func TestSpreadSheetSortByColumnByName(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				spreadSheet := NewSpreadSheet()

				err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
				require.NoError(t, err)

				err = spreadSheet.AddRow([]string{"z", "hello"})
				require.NoError(t, err)

				err = spreadSheet.AddRow([]string{"a", "world"})
				require.NoError(t, err)

				numberOfColumns, err := spreadSheet.GetNumberOfColumns()
				require.NoError(t, err)
				require.EqualValues(t, 2, numberOfColumns)

				numberOfRows, err := spreadSheet.GetNumberOfRows()
				require.NoError(t, err)
				require.EqualValues(t, 2, numberOfRows)

				cell, err := spreadSheet.GetCellValueAsString(0, 0)
				require.NoError(t, err)
				require.EqualValues(t, "z", cell)

				cell, err = spreadSheet.GetCellValueAsString(0, 1)
				require.NoError(t, err)
				require.EqualValues(t, "hello", cell)

				cell, err = spreadSheet.GetCellValueAsString(1, 0)
				require.NoError(t, err)
				require.EqualValues(t, "a", cell)

				cell, err = spreadSheet.GetCellValueAsString(1, 1)
				require.NoError(t, err)
				require.EqualValues(t, "world", cell)

				for i := 0; i < 2; i++ {
					spreadSheet.SortByColumnByName("title1")

					cell, err = spreadSheet.GetCellValueAsString(0, 0)
					require.NoError(t, err)
					require.EqualValues(t, "a", cell)

					cell, err = spreadSheet.GetCellValueAsString(0, 1)
					require.NoError(t, err)
					require.EqualValues(t, "world", cell)

					cell, err = spreadSheet.GetCellValueAsString(1, 0)
					require.NoError(t, err)
					require.EqualValues(t, "z", cell)

					cell, err = spreadSheet.GetCellValueAsString(1, 1)
					require.NoError(t, err)
					require.EqualValues(t, "hello", cell)
				}

				for i := 0; i < 2; i++ {
					spreadSheet.SortByColumnByName("title2")

					cell, err = spreadSheet.GetCellValueAsString(0, 0)
					require.NoError(t, err)
					require.EqualValues(t, "z", cell)

					cell, err = spreadSheet.GetCellValueAsString(0, 1)
					require.NoError(t, err)
					require.EqualValues(t, "hello", cell)

					cell, err = spreadSheet.GetCellValueAsString(1, 0)
					require.NoError(t, err)
					require.EqualValues(t, "a", cell)

					cell, err = spreadSheet.GetCellValueAsString(1, 1)
					require.NoError(t, err)
					require.EqualValues(t, "world", cell)
				}
			},
		)
	}
}

func TestSpreadSheetRenderStringWithoutTitleAndDelimiter(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				spreadSheet := NewSpreadSheet()

				err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
				require.NoError(t, err)

				err = spreadSheet.AddRow([]string{"z", "hello"})
				require.NoError(t, err)

				err = spreadSheet.AddRow([]string{"a", "world"})
				require.NoError(t, err)

				numberOfColumns, err := spreadSheet.GetNumberOfColumns()
				require.NoError(t, err)
				require.EqualValues(t, 2, numberOfColumns)

				numberOfRows, err := spreadSheet.GetNumberOfRows()
				require.NoError(t, err)
				require.EqualValues(t, 2, numberOfRows)

				rendered, err := spreadSheet.RenderAsString(
					&SpreadSheetRenderOptions{
						SkipTitle: true,
						Verbose:   verbose,
					},
				)
				require.NoError(t, err)

				expectedRendered := "z hello\na world\n"
				require.EqualValues(t, expectedRendered, rendered)
			},
		)
	}
}

func TestSpreadSheetRenderStringWithTitleAndWithoutDelimiter(t *testing.T) {
	t.Run("without prefix and suffix", func(t *testing.T) {
		const verbose bool = true

		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"z", "hello"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"a", "world"})
		require.NoError(t, err)

		numberOfColumns, err := spreadSheet.GetNumberOfColumns()
		require.NoError(t, err)
		require.EqualValues(t, 2, numberOfColumns)

		numberOfRows, err := spreadSheet.GetNumberOfRows()
		require.NoError(t, err)
		require.EqualValues(t, 2, numberOfRows)

		rendered, err := spreadSheet.RenderAsString(
			&SpreadSheetRenderOptions{
				SkipTitle: false,
				Verbose:   verbose,
			},
		)
		require.NoError(t, err)

		expectedRendered := "title1 title2\nz hello\na world\n"
		require.EqualValues(t, expectedRendered, rendered)
	})

	t.Run("with prefix and suffix", func(t *testing.T) {
		const verbose bool = true

		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"z", "hello"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"a", "world"})
		require.NoError(t, err)

		numberOfColumns, err := spreadSheet.GetNumberOfColumns()
		require.NoError(t, err)
		require.EqualValues(t, 2, numberOfColumns)

		numberOfRows, err := spreadSheet.GetNumberOfRows()
		require.NoError(t, err)
		require.EqualValues(t, 2, numberOfRows)

		rendered, err := spreadSheet.RenderAsString(
			&SpreadSheetRenderOptions{
				SkipTitle: false,
				Verbose:   verbose,
				Prefix:    "|",
				Suffix:    "|",
			},
		)
		require.NoError(t, err)

		expectedRendered := "| title1 title2 |\n| z hello |\n| a world |\n"
		require.EqualValues(t, expectedRendered, rendered)
	})
}

func TestSpreadSheetRenderStringWithTitleAndDelimiter(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				spreadSheet := NewSpreadSheet()

				err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
				require.NoError(t, err)

				err = spreadSheet.AddRow([]string{"z", "hello"})
				require.NoError(t, err)

				err = spreadSheet.AddRow([]string{"a", "world"})
				require.NoError(t, err)

				numberOfColumns, err := spreadSheet.GetNumberOfColumns()
				require.NoError(t, err)
				require.EqualValues(t, 2, numberOfColumns)

				numberOfRows, err := spreadSheet.GetNumberOfRows()
				require.NoError(t, err)
				require.EqualValues(t, 2, numberOfRows)

				rendered, err := spreadSheet.RenderAsString(
					&SpreadSheetRenderOptions{
						SkipTitle:       false,
						Verbose:         verbose,
						StringDelimiter: "|",
					},
				)
				require.NoError(t, err)

				expectedRendered := "title1 | title2\nz | hello\na | world\n"
				require.EqualValues(t, expectedRendered, rendered)
			},
		)
	}
}

func TestSpreadSheetRenderStringOnlyTitle(t *testing.T) {
	t.Run("only title", func(t *testing.T) {
		const verbose bool = true

		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		numberOfColumns, err := spreadSheet.GetNumberOfColumns()
		require.NoError(t, err)
		require.EqualValues(t, 2, numberOfColumns)

		rendered, err := spreadSheet.RenderAsString(
			&SpreadSheetRenderOptions{
				SkipTitle:       false,
				Verbose:         verbose,
				StringDelimiter: "|",
			},
		)
		require.NoError(t, err)

		expectedRendered := "title1 | title2\n"
		require.EqualValues(t, expectedRendered, rendered)
	})

	t.Run("only title SameColumnWidthForAllRows", func(t *testing.T) {
		const verbose bool = true

		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		numberOfColumns, err := spreadSheet.GetNumberOfColumns()
		require.NoError(t, err)
		require.EqualValues(t, 2, numberOfColumns)

		rendered, err := spreadSheet.RenderAsString(
			&SpreadSheetRenderOptions{
				SkipTitle:                 false,
				Verbose:                   verbose,
				SameColumnWidthForAllRows: true,
				StringDelimiter:           "|",
			},
		)
		require.NoError(t, err)

		expectedRendered := "title1 | title2\n"
		require.EqualValues(t, expectedRendered, rendered)
	})

	t.Run("only title with left and right marker", func(t *testing.T) {
		const verbose bool = true

		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		numberOfColumns, err := spreadSheet.GetNumberOfColumns()
		require.NoError(t, err)
		require.EqualValues(t, 2, numberOfColumns)

		rendered, err := spreadSheet.RenderAsString(
			&SpreadSheetRenderOptions{
				SkipTitle:       false,
				Verbose:         verbose,
				StringDelimiter: "|",
				Prefix:          "|",
				Suffix:          "|",
			},
		)
		require.NoError(t, err)

		expectedRendered := "| title1 | title2 |\n"
		require.EqualValues(t, expectedRendered, rendered)
	})

	t.Run("only title with left and right marker and underline title without cross", func(t *testing.T) {
		const verbose bool = true

		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		numberOfColumns, err := spreadSheet.GetNumberOfColumns()
		require.NoError(t, err)
		require.EqualValues(t, 2, numberOfColumns)

		rendered, err := spreadSheet.RenderAsString(
			&SpreadSheetRenderOptions{
				SkipTitle:       false,
				TitleUnderline:  "-",
				Verbose:         verbose,
				StringDelimiter: "|",
				Prefix:          "|",
				Suffix:          "|",
			},
		)
		require.NoError(t, err)

		expectedRendered := "| title1 | title2 |\n| ------ | ------ |\n"
		require.EqualValues(t, expectedRendered, rendered)
	})
}

func TestSpreadSheetRemoveColumnByName_title1(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				spreadSheet := NewSpreadSheet()

				err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
				require.NoError(t, err)

				err = spreadSheet.AddRow([]string{"z", "hello"})
				require.NoError(t, err)

				err = spreadSheet.AddRow([]string{"a", "world"})
				require.NoError(t, err)

				numberOfColumns, err := spreadSheet.GetNumberOfColumns()
				require.NoError(t, err)
				require.EqualValues(t, 2, numberOfColumns)

				numberOfRows, err := spreadSheet.GetNumberOfRows()
				require.NoError(t, err)
				require.EqualValues(t, 2, numberOfRows)

				err = spreadSheet.RemoveColumnByName("title1")
				require.NoError(t, err)

				numberOfColumns, err = spreadSheet.GetNumberOfColumns()
				require.NoError(t, err)
				require.EqualValues(t, 1, numberOfColumns)

				numberOfRows, err = spreadSheet.GetNumberOfRows()
				require.NoError(t, err)
				require.EqualValues(t, 2, numberOfRows)

				columnTitle, err := spreadSheet.GetColumnTitleAtIndexAsString(0)
				require.NoError(t, err)
				require.EqualValues(t, "title2", columnTitle)

				cell, err := spreadSheet.GetCellValueAsString(0, 0)
				require.NoError(t, err)
				require.EqualValues(t, "hello", cell)

				cell, err = spreadSheet.GetCellValueAsString(1, 0)
				require.NoError(t, err)
				require.EqualValues(t, "world", cell)
			},
		)
	}
}

func TestSpreadSheetRemoveColumnByName_title2(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				spreadSheet := NewSpreadSheet()

				err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
				require.NoError(t, err)

				err = spreadSheet.AddRow([]string{"z", "hello"})
				require.NoError(t, err)

				err = spreadSheet.AddRow([]string{"a", "world"})
				require.NoError(t, err)

				numberOfColumns, err := spreadSheet.GetNumberOfColumns()
				require.NoError(t, err)
				require.EqualValues(t, 2, numberOfColumns)

				numberOfRows, err := spreadSheet.GetNumberOfRows()
				require.NoError(t, err)
				require.EqualValues(t, 2, numberOfRows)

				err = spreadSheet.RemoveColumnByName("title2")
				require.NoError(t, err)

				numberOfColumns, err = spreadSheet.GetNumberOfColumns()
				require.NoError(t, err)
				require.EqualValues(t, 1, numberOfColumns)

				numberOfRows, err = spreadSheet.GetNumberOfRows()
				require.NoError(t, err)
				require.EqualValues(t, 2, numberOfRows)

				columnTitle, err := spreadSheet.GetColumnTitleAtIndexAsString(0)
				require.NoError(t, err)
				require.EqualValues(t, "title1", columnTitle)

				cell, err := spreadSheet.GetCellValueAsString(0, 0)
				require.NoError(t, err)
				require.EqualValues(t, "z", cell)

				cell, err = spreadSheet.GetCellValueAsString(1, 0)
				require.NoError(t, err)
				require.EqualValues(t, "a", cell)
			},
		)
	}
}

func Test_GetRowByFirstColumnValue(t *testing.T) {
	t.Run("empty table", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()
		_, err := spreadSheet.GetRowByFirstColumnValue("a")
		require.Error(t, err)
	})

	t.Run("no rows", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		_, err = spreadSheet.GetRowByFirstColumnValue("a")
		require.Error(t, err)
	})

	t.Run("Entry not found", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"b", "c"})
		require.NoError(t, err)

		_, err = spreadSheet.GetRowByFirstColumnValue("a")
		require.Error(t, err)
	})

	t.Run("Entry found", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"b", "c"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"a", "d"})
		require.NoError(t, err)

		row, err := spreadSheet.GetRowByFirstColumnValue("a")
		require.NoError(t, err)

		entries, err := row.GetEntries()
		require.NoError(t, err)
		require.EqualValues(t, []string{"a", "d"}, entries)
	})

	t.Run("Empty first column value found", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"b", "c"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"a", "d"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"", "only second cell contains value"})
		require.NoError(t, err)

		row, err := spreadSheet.GetRowByFirstColumnValue("")
		require.NoError(t, err)

		entries, err := row.GetEntries()
		require.NoError(t, err)
		require.EqualValues(t, []string{"", "only second cell contains value"}, entries)
	})
}

func Test_UpdateRowFoundByFirstColumnValue(t *testing.T) {
	t.Run("empty table", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()
		err := spreadSheet.UpdateRowFoundByFirstColumnValue("a", 1, "updated value")
		require.Error(t, err)
	})

	t.Run("no rows", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		err = spreadSheet.UpdateRowFoundByFirstColumnValue("a", 1, "updated value")
		require.Error(t, err)
	})

	t.Run("Entry not found", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"b", "c"})
		require.NoError(t, err)

		err = spreadSheet.UpdateRowFoundByFirstColumnValue("a", 1, "updated value")
		require.Error(t, err)
	})

	t.Run("Invalid index", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"b", "c"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"a", "b"})
		require.NoError(t, err)

		err = spreadSheet.UpdateRowFoundByFirstColumnValue("a", -1, "updated value")
		require.Error(t, err)
	})

	t.Run("Update first cell", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"b", "c"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"a", "b"})
		require.NoError(t, err)

		err = spreadSheet.UpdateRowFoundByFirstColumnValue("a", 0, "updated value")
		require.NoError(t, err)

		row0, err := spreadSheet.GetRowByIndexAsStringSlice(0)
		require.NoError(t, err)
		require.EqualValues(t, []string{"b", "c"}, row0)

		row1, err := spreadSheet.GetRowByIndexAsStringSlice(1)
		require.NoError(t, err)
		require.EqualValues(t, []string{"updated value", "b"}, row1)
	})

	t.Run("Update second cell", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"b", "c"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"a", "b"})
		require.NoError(t, err)

		err = spreadSheet.UpdateRowFoundByFirstColumnValue("a", 1, "updated value")
		require.NoError(t, err)

		row0, err := spreadSheet.GetRowByIndexAsStringSlice(0)
		require.NoError(t, err)
		require.EqualValues(t, []string{"b", "c"}, row0)

		row1, err := spreadSheet.GetRowByIndexAsStringSlice(1)
		require.NoError(t, err)
		require.EqualValues(t, []string{"a", "updated value"}, row1)
	})

	t.Run("Too hight index", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()

		err := spreadSheet.SetColumnTitles([]string{"title1", "title2"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"b", "c"})
		require.NoError(t, err)

		err = spreadSheet.AddRow([]string{"a", "b"})
		require.NoError(t, err)

		err = spreadSheet.UpdateRowFoundByFirstColumnValue("a", 2, "updated value")
		require.Error(t, err)
	})
}

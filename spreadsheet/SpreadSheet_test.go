package spreadsheet

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
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

				require.EqualValues(t, 0, spreadSheet.MustGetNumberOfColumns())
				require.EqualValues(t, 0, spreadSheet.MustGetNumberOfRows())
				require.True(t, spreadSheet.MustIsEmpty())
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
				require := require.New(t)

				spreadSheet := NewSpreadSheet()
				spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})

				require.EqualValues(2, spreadSheet.MustGetNumberOfColumns())
				require.EqualValues(0, spreadSheet.MustGetNumberOfRows())
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
				require := require.New(t)

				spreadSheet := NewSpreadSheet()
				spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
				spreadSheet.MustAddRow([]string{"z", "hello"})
				spreadSheet.MustAddRow([]string{"a", "world"})

				require.EqualValues(2, spreadSheet.MustGetNumberOfColumns())
				require.EqualValues(2, spreadSheet.MustGetNumberOfRows())

				require.EqualValues("z", spreadSheet.MustGetCellValueAsString(0, 0))
				require.EqualValues("hello", spreadSheet.MustGetCellValueAsString(0, 1))
				require.EqualValues("a", spreadSheet.MustGetCellValueAsString(1, 0))
				require.EqualValues("world", spreadSheet.MustGetCellValueAsString(1, 1))

				for i := 0; i < 2; i++ {
					spreadSheet.SortByColumnByName("title1")
					require.EqualValues("a", spreadSheet.MustGetCellValueAsString(0, 0))
					require.EqualValues("world", spreadSheet.MustGetCellValueAsString(0, 1))
					require.EqualValues("z", spreadSheet.MustGetCellValueAsString(1, 0))
					require.EqualValues("hello", spreadSheet.MustGetCellValueAsString(1, 1))
				}

				for i := 0; i < 2; i++ {
					spreadSheet.SortByColumnByName("title2")
					require.EqualValues("z", spreadSheet.MustGetCellValueAsString(0, 0))
					require.EqualValues("hello", spreadSheet.MustGetCellValueAsString(0, 1))
					require.EqualValues("a", spreadSheet.MustGetCellValueAsString(1, 0))
					require.EqualValues("world", spreadSheet.MustGetCellValueAsString(1, 1))
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
				require := require.New(t)

				const verbose bool = true

				spreadSheet := NewSpreadSheet()
				spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
				spreadSheet.MustAddRow([]string{"z", "hello"})
				spreadSheet.MustAddRow([]string{"a", "world"})

				require.EqualValues(2, spreadSheet.MustGetNumberOfColumns())
				require.EqualValues(2, spreadSheet.MustGetNumberOfRows())

				rendered := spreadSheet.MustRenderAsString(
					&SpreadSheetRenderOptions{
						SkipTitle: true,
						Verbose:   verbose,
					},
				)

				expectedRendered := "z hello\na world\n"
				require.EqualValues(expectedRendered, rendered)
			},
		)
	}
}

func TestSpreadSheetRenderStringWithTitleAndWithoutDelimiter(t *testing.T) {
	t.Run("without prefix and suffix", func(t *testing.T) {
		require := require.New(t)

		const verbose bool = true

		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
		spreadSheet.MustAddRow([]string{"z", "hello"})
		spreadSheet.MustAddRow([]string{"a", "world"})

		require.EqualValues(2, spreadSheet.MustGetNumberOfColumns())
		require.EqualValues(2, spreadSheet.MustGetNumberOfRows())

		rendered := spreadSheet.MustRenderAsString(
			&SpreadSheetRenderOptions{
				SkipTitle: false,
				Verbose:   verbose,
			},
		)

		expectedRendered := "title1 title2\nz hello\na world\n"
		require.EqualValues(expectedRendered, rendered)
	})

	t.Run("with prefix and suffix", func(t *testing.T) {
		require := require.New(t)

		const verbose bool = true

		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
		spreadSheet.MustAddRow([]string{"z", "hello"})
		spreadSheet.MustAddRow([]string{"a", "world"})

		require.EqualValues(2, spreadSheet.MustGetNumberOfColumns())
		require.EqualValues(2, spreadSheet.MustGetNumberOfRows())

		rendered := spreadSheet.MustRenderAsString(
			&SpreadSheetRenderOptions{
				SkipTitle: false,
				Verbose:   verbose,
				Prefix:    "|",
				Suffix:    "|",
			},
		)

		expectedRendered := "| title1 title2 |\n| z hello |\n| a world |\n"
		require.EqualValues(expectedRendered, rendered)
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
				require := require.New(t)

				const verbose bool = true

				spreadSheet := NewSpreadSheet()
				spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
				spreadSheet.MustAddRow([]string{"z", "hello"})
				spreadSheet.MustAddRow([]string{"a", "world"})

				require.EqualValues(2, spreadSheet.MustGetNumberOfColumns())
				require.EqualValues(2, spreadSheet.MustGetNumberOfRows())

				rendered := spreadSheet.MustRenderAsString(
					&SpreadSheetRenderOptions{
						SkipTitle:       false,
						Verbose:         verbose,
						StringDelimiter: "|",
					},
				)

				expectedRendered := "title1 | title2\nz | hello\na | world\n"
				require.EqualValues(expectedRendered, rendered)
			},
		)
	}
}

func TestSpreadSheetRenderStringOnlyTitle(t *testing.T) {
	t.Run("only title", func(t *testing.T) {
		require := require.New(t)

		const verbose bool = true

		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})

		require.EqualValues(2, spreadSheet.MustGetNumberOfColumns())

		rendered := spreadSheet.MustRenderAsString(
			&SpreadSheetRenderOptions{
				SkipTitle:       false,
				Verbose:         verbose,
				StringDelimiter: "|",
			},
		)

		expectedRendered := "title1 | title2\n"
		require.EqualValues(expectedRendered, rendered)
	})

	t.Run("only title SameColumnWidthForAllRows", func(t *testing.T) {
		require := require.New(t)

		const verbose bool = true

		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})

		require.EqualValues(2, spreadSheet.MustGetNumberOfColumns())

		rendered := spreadSheet.MustRenderAsString(
			&SpreadSheetRenderOptions{
				SkipTitle:                 false,
				Verbose:                   verbose,
				SameColumnWidthForAllRows: true,
				StringDelimiter:           "|",
			},
		)

		expectedRendered := "title1 | title2\n"
		require.EqualValues(expectedRendered, rendered)
	})

	t.Run("only title with left and right marker", func(t *testing.T) {
		require := require.New(t)

		const verbose bool = true

		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})

		require.EqualValues(2, spreadSheet.MustGetNumberOfColumns())

		rendered := spreadSheet.MustRenderAsString(
			&SpreadSheetRenderOptions{
				SkipTitle:       false,
				Verbose:         verbose,
				StringDelimiter: "|",
				Prefix:          "|",
				Suffix:          "|",
			},
		)

		expectedRendered := "| title1 | title2 |\n"
		require.EqualValues(expectedRendered, rendered)
	})

	t.Run("only title with left and right marker and underline title without cross", func(t *testing.T) {
		require := require.New(t)

		const verbose bool = true

		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})

		require.EqualValues(2, spreadSheet.MustGetNumberOfColumns())

		rendered := spreadSheet.MustRenderAsString(
			&SpreadSheetRenderOptions{
				SkipTitle:       false,
				TitleUnderline:  "-",
				Verbose:         verbose,
				StringDelimiter: "|",
				Prefix:          "|",
				Suffix:          "|",
			},
		)

		expectedRendered := "| title1 | title2 |\n| ------ | ------ |\n"
		require.EqualValues(expectedRendered, rendered)
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
				require := require.New(t)

				spreadSheet := NewSpreadSheet()
				spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
				spreadSheet.MustAddRow([]string{"z", "hello"})
				spreadSheet.MustAddRow([]string{"a", "world"})

				require.EqualValues(2, spreadSheet.MustGetNumberOfColumns())
				require.EqualValues(2, spreadSheet.MustGetNumberOfRows())

				spreadSheet.MustRemoveColumnByName("title1")
				require.EqualValues(1, spreadSheet.MustGetNumberOfColumns())
				require.EqualValues(2, spreadSheet.MustGetNumberOfRows())

				require.EqualValues("title2", spreadSheet.MustGetColumnTitleAtIndexAsString(0))
				require.EqualValues("hello", spreadSheet.MustGetCellValueAsString(0, 0))
				require.EqualValues("world", spreadSheet.MustGetCellValueAsString(1, 0))
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
				require := require.New(t)

				spreadSheet := NewSpreadSheet()
				spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
				spreadSheet.MustAddRow([]string{"z", "hello"})
				spreadSheet.MustAddRow([]string{"a", "world"})

				require.EqualValues(2, spreadSheet.MustGetNumberOfColumns())
				require.EqualValues(2, spreadSheet.MustGetNumberOfRows())

				spreadSheet.MustRemoveColumnByName("title2")
				require.EqualValues(1, spreadSheet.MustGetNumberOfColumns())
				require.EqualValues(2, spreadSheet.MustGetNumberOfRows())

				require.EqualValues("title1", spreadSheet.MustGetColumnTitleAtIndexAsString(0))
				require.EqualValues("z", spreadSheet.MustGetCellValueAsString(0, 0))
				require.EqualValues("a", spreadSheet.MustGetCellValueAsString(1, 0))
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
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})

		_, err := spreadSheet.GetRowByFirstColumnValue("a")
		require.Error(t, err)
	})

	t.Run("Entry not found", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
		spreadSheet.MustAddRow([]string{"b", "c"})

		_, err := spreadSheet.GetRowByFirstColumnValue("a")
		require.Error(t, err)
	})

	t.Run("Entry found", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
		spreadSheet.MustAddRow([]string{"b", "c"})
		spreadSheet.MustAddRow([]string{"a", "d"})

		row, err := spreadSheet.GetRowByFirstColumnValue("a")
		require.NoError(t, err)

		entries, err := row.GetEntries()
		require.NoError(t, err)
		require.EqualValues(t, []string{"a", "d"}, entries)
	})

	t.Run("Entry found", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
		spreadSheet.MustAddRow([]string{"b", "c"})
		spreadSheet.MustAddRow([]string{"a", "d"})
		spreadSheet.MustAddRow([]string{"", "only second cell contains value"})

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
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})

		err := spreadSheet.UpdateRowFoundByFirstColumnValue("a", 1, "updated value")
		require.Error(t, err)
	})

	t.Run("Entry not found", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
		spreadSheet.MustAddRow([]string{"b", "c"})

		err := spreadSheet.UpdateRowFoundByFirstColumnValue("a", 1, "updated value")
		require.Error(t, err)
	})

	t.Run("Invalid index", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
		spreadSheet.MustAddRow([]string{"b", "c"})
		spreadSheet.MustAddRow([]string{"a", "b"})

		err := spreadSheet.UpdateRowFoundByFirstColumnValue("a", -1, "updated value")
		require.Error(t, err)
	})

	t.Run("Update first cell", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
		spreadSheet.MustAddRow([]string{"b", "c"})
		spreadSheet.MustAddRow([]string{"a", "b"})

		err := spreadSheet.UpdateRowFoundByFirstColumnValue("a", 0, "updated value")
		require.NoError(t, err)

		require.EqualValues(
			t,
			[]string{"b", "c"},
			spreadSheet.MustGetRowByIndexAsStringSlice(0),
		)

		require.EqualValues(
			t,
			[]string{"updated value", "b"},
			spreadSheet.MustGetRowByIndexAsStringSlice(1),
		)
	})

	t.Run("Update second cell", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
		spreadSheet.MustAddRow([]string{"b", "c"})
		spreadSheet.MustAddRow([]string{"a", "b"})

		err := spreadSheet.UpdateRowFoundByFirstColumnValue("a", 1, "updated value")
		require.NoError(t, err)

		require.EqualValues(
			t,
			[]string{"b", "c"},
			spreadSheet.MustGetRowByIndexAsStringSlice(0),
		)

		require.EqualValues(
			t,
			[]string{"a", "updated value"},
			spreadSheet.MustGetRowByIndexAsStringSlice(1),
		)
	})

	t.Run("Too hight index", func(t *testing.T) {
		spreadSheet := NewSpreadSheet()
		spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
		spreadSheet.MustAddRow([]string{"b", "c"})
		spreadSheet.MustAddRow([]string{"a", "b"})

		err := spreadSheet.UpdateRowFoundByFirstColumnValue("a", 2, "updated value")
		require.Error(t, err)
	})
}

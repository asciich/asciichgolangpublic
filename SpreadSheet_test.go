package asciichgolangpublic

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
				require := require.New(t)

				spreadSheet := NewSpreadSheet()

				require.EqualValues(0, spreadSheet.MustGetNumberOfColumns())
				require.EqualValues(0, spreadSheet.MustGetNumberOfRows())
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

package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert := assert.New(t)

				spreadSheet := NewSpreadSheet()

				assert.EqualValues(0, spreadSheet.MustGetNumberOfColumns())
				assert.EqualValues(0, spreadSheet.MustGetNumberOfRows())
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
				assert := assert.New(t)

				spreadSheet := NewSpreadSheet()
				spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})

				assert.EqualValues(2, spreadSheet.MustGetNumberOfColumns())
				assert.EqualValues(0, spreadSheet.MustGetNumberOfRows())
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
				assert := assert.New(t)

				spreadSheet := NewSpreadSheet()
				spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
				spreadSheet.MustAddRow([]string{"z", "hello"})
				spreadSheet.MustAddRow([]string{"a", "world"})

				assert.EqualValues(2, spreadSheet.MustGetNumberOfColumns())
				assert.EqualValues(2, spreadSheet.MustGetNumberOfRows())

				assert.EqualValues("z", spreadSheet.MustGetCellValueAsString(0, 0))
				assert.EqualValues("hello", spreadSheet.MustGetCellValueAsString(0, 1))
				assert.EqualValues("a", spreadSheet.MustGetCellValueAsString(1, 0))
				assert.EqualValues("world", spreadSheet.MustGetCellValueAsString(1, 1))

				for i := 0; i < 2; i++ {
					spreadSheet.SortByColumnByName("title1")
					assert.EqualValues("a", spreadSheet.MustGetCellValueAsString(0, 0))
					assert.EqualValues("world", spreadSheet.MustGetCellValueAsString(0, 1))
					assert.EqualValues("z", spreadSheet.MustGetCellValueAsString(1, 0))
					assert.EqualValues("hello", spreadSheet.MustGetCellValueAsString(1, 1))
				}

				for i := 0; i < 2; i++ {
					spreadSheet.SortByColumnByName("title2")
					assert.EqualValues("z", spreadSheet.MustGetCellValueAsString(0, 0))
					assert.EqualValues("hello", spreadSheet.MustGetCellValueAsString(0, 1))
					assert.EqualValues("a", spreadSheet.MustGetCellValueAsString(1, 0))
					assert.EqualValues("world", spreadSheet.MustGetCellValueAsString(1, 1))
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
				assert := assert.New(t)

				const verbose bool = true

				spreadSheet := NewSpreadSheet()
				spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
				spreadSheet.MustAddRow([]string{"z", "hello"})
				spreadSheet.MustAddRow([]string{"a", "world"})

				assert.EqualValues(2, spreadSheet.MustGetNumberOfColumns())
				assert.EqualValues(2, spreadSheet.MustGetNumberOfRows())

				rendered := spreadSheet.MustRenderAsString(
					&SpreadSheetRenderOptions{
						SkipTitle: true,
						Verbose:   verbose,
					},
				)

				expectedRendered := "z hello\na world\n"
				assert.EqualValues(expectedRendered, rendered)
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
				assert := assert.New(t)

				spreadSheet := NewSpreadSheet()
				spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
				spreadSheet.MustAddRow([]string{"z", "hello"})
				spreadSheet.MustAddRow([]string{"a", "world"})

				assert.EqualValues(2, spreadSheet.MustGetNumberOfColumns())
				assert.EqualValues(2, spreadSheet.MustGetNumberOfRows())

				spreadSheet.MustRemoveColumnByName("title1")
				assert.EqualValues(1, spreadSheet.MustGetNumberOfColumns())
				assert.EqualValues(2, spreadSheet.MustGetNumberOfRows())

				assert.EqualValues("title2", spreadSheet.MustGetColumnTitleAtIndexAsString(0))
				assert.EqualValues("hello", spreadSheet.MustGetCellValueAsString(0, 0))
				assert.EqualValues("world", spreadSheet.MustGetCellValueAsString(1, 0))
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
				assert := assert.New(t)

				spreadSheet := NewSpreadSheet()
				spreadSheet.MustSetColumnTitles([]string{"title1", "title2"})
				spreadSheet.MustAddRow([]string{"z", "hello"})
				spreadSheet.MustAddRow([]string{"a", "world"})

				assert.EqualValues(2, spreadSheet.MustGetNumberOfColumns())
				assert.EqualValues(2, spreadSheet.MustGetNumberOfRows())

				spreadSheet.MustRemoveColumnByName("title2")
				assert.EqualValues(1, spreadSheet.MustGetNumberOfColumns())
				assert.EqualValues(2, spreadSheet.MustGetNumberOfRows())

				assert.EqualValues("title1", spreadSheet.MustGetColumnTitleAtIndexAsString(0))
				assert.EqualValues("z", spreadSheet.MustGetCellValueAsString(0, 0))
				assert.EqualValues("a", spreadSheet.MustGetCellValueAsString(1, 0))
			},
		)
	}
}

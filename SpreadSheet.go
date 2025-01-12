package asciichgolangpublic

import (
	"fmt"
	"sort"
	"strconv"

	astrings "github.com/asciich/asciichgolangpublic/datatypes/strings"
)

type SpreadSheet struct {
	TitleRow *SpreadSheetRow
	rows     []*SpreadSheetRow
}

func GetSpreadsheetWithNColumns(nColumns int) (s *SpreadSheet, err error) {
	if nColumns <= 0 {
		return nil, TracedErrorf("Invalid nColumns: '%d'", nColumns)
	}

	titles := []string{}
	for i := 0; i < nColumns; i++ {
		titles = append(titles, "column"+strconv.Itoa(i))
	}

	s = NewSpreadSheet()
	err = s.SetColumnTitles(titles)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func MustGetSpreadsheetWithNColumns(nColumns int) (s *SpreadSheet) {
	s, err := GetSpreadsheetWithNColumns(nColumns)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return s
}

func NewSpreadSheet() (s *SpreadSheet) {
	return new(SpreadSheet)
}

func (s *SpreadSheet) AddRow(rowEntries []string) (err error) {
	if len(rowEntries) <= 0 {
		return TracedError("rowEntries is empty")
	}

	nColumns, err := s.GetNumberOfColumns()
	if err != nil {
		return err
	}

	nColumnsToAdd := len(rowEntries)
	if nColumnsToAdd != nColumns {
		return TracedErrorf(
			"Number of columns mismatch. Row to add has '%d' columns but spreadsheet has '%d' columns.",
			nColumnsToAdd,
			nColumns,
		)
	}

	rowToAdd := NewSpreadSheetRow()
	err = rowToAdd.SetEntries(rowEntries)
	if err != nil {
		return err
	}

	s.rows = append(s.rows, rowToAdd)

	return nil
}

func (s *SpreadSheet) GetCellValueAsString(rowIndex int, columnIndex int) (cellValue string, err error) {
	if rowIndex < 0 {
		return "", TracedErrorf("Invalid rowIndex: '%d'", rowIndex)
	}

	if columnIndex < 0 {
		return "", TracedErrorf("Invalid columnIndex: '%d'", columnIndex)
	}

	row, err := s.GetRowByIndex(rowIndex)
	if err != nil {
		return "", err
	}

	cellValue, err = row.GetColumnValueAsString(columnIndex)
	if err != nil {
		return "", err
	}

	return cellValue, nil
}

func (s *SpreadSheet) GetColumnIndexByName(columnName string) (columnIndex int, err error) {
	if columnName == "" {
		return -1, TracedError("columnName is empty string")
	}

	titles, err := s.GetColumnTitlesAsStringSlice()
	if err != nil {
		return -1, err
	}

	for i, title := range titles {
		if title == columnName {
			return i, nil
		}
	}

	return -1, TracedErrorf("Unable to find column title '%s'", columnName)
}

func (s *SpreadSheet) GetColumnTitleAtIndexAsString(index int) (title string, err error) {
	titleRow, err := s.GetTitleRow()
	if err != nil {
		return "", err
	}

	title, err = titleRow.GetColumnValueAsString(index)
	if err != nil {
		return "", err
	}

	return title, nil
}

func (s *SpreadSheet) GetColumnTitlesAsStringSlice() (titles []string, err error) {
	TitleRow, err := s.GetTitleRow()
	if err != nil {
		return nil, err
	}

	titles, err = TitleRow.GetEntries()
	if err != nil {
		return nil, err
	}

	return titles, nil
}

func (s *SpreadSheet) GetMaxColumnWidths() (columnWitdhs []int, err error) {
	nColumns, err := s.GetNumberOfColumns()
	if err != nil {
		return nil, err
	}

	if nColumns <= 0 {
		return []int{}, nil
	}

	rows, err := s.GetRows()
	if err != nil {
		return nil, err
	}

	columnWidths := Slices().GetIntSliceInitializedWithZeros(nColumns)

	for _, row := range rows {
		rowColumnWidths, err := row.GetColumnWidths()
		if err != nil {
			return nil, err
		}

		columnWidths = Slices().MaxIntValuePerIndex(columnWidths, rowColumnWidths)
	}

	return columnWidths, nil
}

func (s *SpreadSheet) GetMinColumnWithsAsSelectedInOptions(options *SpreadSheetRenderOptions) (columnWidths []int, err error) {
	if options == nil {
		return nil, TracedError("options is nil")
	}

	nColumns, err := s.GetNumberOfColumns()
	if err != nil {
		return nil, err
	}

	if !options.SameColumnWidthForAllRows {
		columnWidths = []int{}
		for i := 0; i < nColumns; i++ {
			columnWidths = append(columnWidths, 0)
		}

		return columnWidths, nil
	}

	columnWidths, err = s.GetMaxColumnWidths()
	if err != nil {
		return nil, err
	}

	return columnWidths, nil
}

func (s *SpreadSheet) GetNumberOfColumns() (nColumns int, err error) {
	if s.TitleRow == nil {
		return 0, nil
	}

	TitleRow, err := s.GetTitleRow()
	if err != nil {
		return -1, err
	}

	nColumns, err = TitleRow.GetNumberOfEntries()
	if err != nil {
		return -1, err
	}

	return nColumns, nil
}

func (s *SpreadSheet) GetNumberOfRows() (nRows int, err error) {
	if s.rows == nil {
		return 0, nil
	}

	rows, err := s.GetRows()
	if err != nil {
		return -1, err
	}

	return len(rows), nil
}

func (s *SpreadSheet) GetRowByIndex(rowIndex int) (row *SpreadSheetRow, err error) {
	if rowIndex < 0 {
		return nil, TracedErrorf("Invalid rowIndex: '%d'", rowIndex)
	}

	rows, err := s.GetRows()
	if err != nil {
		return nil, err
	}

	nRows := len(rows)

	if rowIndex >= nRows {
		return nil, TracedErrorf("Invlaid rowIndex: Index '%d' is invalid for a spreadsheet with '%d' rows.", rowIndex, nRows)
	}

	row = rows[rowIndex]

	if row == nil {
		return nil, TracedError("row is nil after evaluation")
	}

	return row, nil
}

func (s *SpreadSheet) GetRows() (rows []*SpreadSheetRow, err error) {
	if s.rows == nil {
		return nil, TracedErrorf("rows not set")
	}

	if len(s.rows) <= 0 {
		return nil, TracedErrorf("rows has no elements")
	}

	return s.rows, nil
}

func (s *SpreadSheet) GetTitleRow() (TitleRow *SpreadSheetRow, err error) {
	if s.TitleRow == nil {
		return nil, TracedErrorf("TitleRow not set")
	}

	return s.TitleRow, nil
}

func (s *SpreadSheet) MustAddRow(rowEntries []string) {
	err := s.AddRow(rowEntries)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustGetCellValueAsString(rowIndex int, columnIndex int) (cellValue string) {
	cellValue, err := s.GetCellValueAsString(rowIndex, columnIndex)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return cellValue
}

func (s *SpreadSheet) MustGetColumnIndexByName(columnName string) (columnIndex int) {
	columnIndex, err := s.GetColumnIndexByName(columnName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return columnIndex
}

func (s *SpreadSheet) MustGetColumnTitleAtIndexAsString(index int) (title string) {
	title, err := s.GetColumnTitleAtIndexAsString(index)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return title
}

func (s *SpreadSheet) MustGetColumnTitlesAsStringSlice() (titles []string) {
	titles, err := s.GetColumnTitlesAsStringSlice()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return titles
}

func (s *SpreadSheet) MustGetMaxColumnWidths() (columnWitdhs []int) {
	columnWitdhs, err := s.GetMaxColumnWidths()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return columnWitdhs
}

func (s *SpreadSheet) MustGetMinColumnWithsAsSelectedInOptions(options *SpreadSheetRenderOptions) (columnWidths []int) {
	columnWidths, err := s.GetMinColumnWithsAsSelectedInOptions(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return columnWidths
}

func (s *SpreadSheet) MustGetNumberOfColumns() (nColumns int) {
	nColumns, err := s.GetNumberOfColumns()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nColumns
}

func (s *SpreadSheet) MustGetNumberOfRows() (nRows int) {
	nRows, err := s.GetNumberOfRows()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nRows
}

func (s *SpreadSheet) MustGetRowByIndex(rowIndex int) (row *SpreadSheetRow) {
	row, err := s.GetRowByIndex(rowIndex)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return row
}

func (s *SpreadSheet) MustGetRows() (rows []*SpreadSheetRow) {
	rows, err := s.GetRows()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return rows
}

func (s *SpreadSheet) MustGetTitleRow() (TitleRow *SpreadSheetRow) {
	TitleRow, err := s.GetTitleRow()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return TitleRow
}

func (s *SpreadSheet) MustPrintAsString(options *SpreadSheetRenderOptions) {
	err := s.PrintAsString(options)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustRemoveColumnByIndex(columnIndex int) {
	err := s.RemoveColumnByIndex(columnIndex)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustRemoveColumnByName(columnName string) {
	err := s.RemoveColumnByName(columnName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustRenderAsString(options *SpreadSheetRenderOptions) (rendered string) {
	rendered, err := s.RenderAsString(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return rendered
}

func (s *SpreadSheet) MustRenderTitleRowAsString(options *SpreadSheetRenderRowOptions) (rendered string) {
	rendered, err := s.RenderTitleRowAsString(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return rendered
}

func (s *SpreadSheet) MustRenderToStdout(options *SpreadSheetRenderOptions) {
	err := s.RenderToStdout(options)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustSetColumnTitles(titles []string) {
	err := s.SetColumnTitles(titles)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustSetRows(rows []*SpreadSheetRow) {
	err := s.SetRows(rows)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustSetTitleRow(TitleRow *SpreadSheetRow) {
	err := s.SetTitleRow(TitleRow)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustSortByColumnByName(columnName string) {
	err := s.SortByColumnByName(columnName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) PrintAsString(options *SpreadSheetRenderOptions) (err error) {
	if options == nil {
		return TracedError("options is nil")
	}

	stringToPrint, err := s.RenderAsString(options)
	if err != nil {
		return err
	}

	fmt.Print(stringToPrint)

	return nil
}

func (s *SpreadSheet) RemoveColumnByIndex(columnIndex int) (err error) {
	if columnIndex < 0 {
		return TracedErrorf("columnIndex '%d' is invalid", columnIndex)
	}

	nColumns, err := s.GetNumberOfColumns()
	if err != nil {
		return err
	}

	if columnIndex >= nColumns {
		return TracedErrorf(
			"Invalid columnIndex '%d' for a spreadsheet with '%d' columns.",
			columnIndex,
			nColumns,
		)
	}

	titleRow, err := s.GetTitleRow()
	if err != nil {
		return err
	}

	err = titleRow.RemoveElementAtIndex(columnIndex)
	if err != nil {
		return err
	}

	rows, err := s.GetRows()
	if err != nil {
		return err
	}

	for _, row := range rows {
		err = row.RemoveElementAtIndex(columnIndex)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SpreadSheet) RemoveColumnByName(columnName string) (err error) {
	if columnName == "" {
		return TracedError("columntName is empty string")
	}

	columnIndex, err := s.GetColumnIndexByName(columnName)
	if err != nil {
		return err
	}

	err = s.RemoveColumnByIndex(columnIndex)
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheet) RenderAsString(options *SpreadSheetRenderOptions) (rendered string, err error) {
	if options == nil {
		return "", TracedErrorNil("options")
	}

	var minColumnWidths []int = nil
	if options.SameColumnWidthForAllRows {
		minColumnWidths, err = s.GetMinColumnWithsAsSelectedInOptions(options)
		if err != nil {
			return "", err
		}
	}

	renderRowOptions := NewSpreadSheetRenderRowOptions()
	renderRowOptions.Verbose = options.Verbose
	renderRowOptions.MinColumnWidths = minColumnWidths
	renderRowOptions.StringDelimiter = options.StringDelimiter

	rendered = ""
	if !options.SkipTitle {
		toAdd, err := s.RenderTitleRowAsString(renderRowOptions)
		if err != nil {
			return "", err
		}

		rendered += toAdd + "\n"
	}

	rows, err := s.GetRows()
	if err != nil {
		return "", err
	}

	for _, row := range rows {
		toAdd, err := row.RenderAsString(renderRowOptions)
		if err != nil {
			return "", err
		}

		rendered += toAdd + "\n"
	}

	return rendered, nil
}

func (s *SpreadSheet) RenderTitleRowAsString(options *SpreadSheetRenderRowOptions) (rendered string, err error) {
	if options == nil {
		return "", TracedError("options is nil")
	}

	titleRow, err := s.GetTitleRow()
	if err != nil {
		return "", err
	}

	rendered, err = titleRow.RenderAsString(options)
	if err != nil {
		return "", err
	}

	return rendered, nil
}

func (s *SpreadSheet) RenderToStdout(options *SpreadSheetRenderOptions) (err error) {
	if options == nil {
		return TracedErrorNil("options")
	}

	rendered, err := s.RenderAsString(options)
	if err != nil {
		return err
	}

	rendered = astrings.EnsureEndsWithLineBreak(rendered)

	fmt.Print(rendered)

	return nil
}

func (s *SpreadSheet) SetColumnTitles(titles []string) (err error) {
	if len(titles) <= 0 {
		return TracedError("titles is empty slice")
	}

	TitleRow := NewSpreadSheetRow()
	err = TitleRow.SetEntries(titles)
	if err != nil {
		return err
	}

	s.TitleRow = TitleRow

	return nil
}

func (s *SpreadSheet) SetRows(rows []*SpreadSheetRow) (err error) {
	if rows == nil {
		return TracedErrorf("rows is nil")
	}

	if len(rows) <= 0 {
		return TracedErrorf("rows has no elements")
	}

	s.rows = rows

	return nil
}

func (s *SpreadSheet) SetTitleRow(TitleRow *SpreadSheetRow) (err error) {
	if TitleRow == nil {
		return TracedErrorf("TitleRow is nil")
	}

	s.TitleRow = TitleRow

	return nil
}

func (s *SpreadSheet) SortByColumnByName(columnName string) (err error) {
	if columnName == "" {
		return TracedError("columnName is empty string")
	}

	columnIndex, err := s.GetColumnIndexByName(columnName)
	if err != nil {
		return err
	}

	rows, err := s.GetRows()
	if err != nil {
		return err
	}

	sort.Slice(rows, func(i int, j int) bool {
		iRow, err := s.GetRowByIndex(i)
		if err != nil {
			LogGoErrorFatal(err)
		}

		jRow, err := s.GetRowByIndex(j)
		if err != nil {
			LogGoErrorFatal(err)
		}

		iValue, err := iRow.GetColumnValueAsString(columnIndex)
		if err != nil {
			LogGoErrorFatal(err)
		}

		jValue, err := jRow.GetColumnValueAsString(columnIndex)
		if err != nil {
			LogGoErrorFatal(err)
		}

		return iValue < jValue
	})

	err = s.SetRows(rows)
	if err != nil {
		return err
	}

	return nil
}

package spreadsheet

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type SpreadSheet struct {
	TitleRow *SpreadSheetRow
	rows     []*SpreadSheetRow
}

func GetSpreadsheetWithNColumns(nColumns int) (s *SpreadSheet, err error) {
	if nColumns <= 0 {
		return nil, tracederrors.TracedErrorf("Invalid nColumns: '%d'", nColumns)
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
		logging.LogGoErrorFatal(err)
	}

	return s
}

func NewSpreadSheet() (s *SpreadSheet) {
	return new(SpreadSheet)
}

func (s *SpreadSheet) MustIsEmpty() (isEmpty bool) {
	isEmpty, err := s.IsEmpty()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isEmpty
}

func (s *SpreadSheet) IsEmpty() (isEmpty bool, err error) {
	ncols, err := s.GetNumberOfColumns()
	if err != nil {
		return false, err
	}

	if ncols > 0 {
		return false, nil
	}

	nrows, err := s.GetNumberOfRows()
	if err != nil {
		return false, err
	}

	if nrows > 0 {
		return false, nil
	}

	return true, nil
}

func (s *SpreadSheet) AddRow(rowEntries []string) (err error) {
	if len(rowEntries) <= 0 {
		return tracederrors.TracedError("rowEntries is empty")
	}

	nColumns, err := s.GetNumberOfColumns()
	if err != nil {
		return err
	}

	nColumnsToAdd := len(rowEntries)
	if nColumnsToAdd != nColumns {
		return tracederrors.TracedErrorf(
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
		return "", tracederrors.TracedErrorf("Invalid rowIndex: '%d'", rowIndex)
	}

	if columnIndex < 0 {
		return "", tracederrors.TracedErrorf("Invalid columnIndex: '%d'", columnIndex)
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
		return -1, tracederrors.TracedError("columnName is empty string")
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

	return -1, tracederrors.TracedErrorf("Unable to find column title '%s'", columnName)
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

	columnWidths := slicesutils.GetIntSliceInitializedWithZeros(nColumns)

	nRows, err := s.GetNumberOfRows()
	if err != nil {
		return nil, err
	}

	if nRows > 0 {
		rows, err := s.GetRows()
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			rowColumnWidths, err := row.GetColumnWidths()
			if err != nil {
				return nil, err
			}

			columnWidths = slicesutils.MaxIntValuePerIndex(columnWidths, rowColumnWidths)
		}
	}

	return columnWidths, nil
}

func (s *SpreadSheet) GetRowByFirstColumnValue(value string) (row *SpreadSheetRow, err error) {
	for _, r := range s.rows {
		if len(r.entries) <= 0 {
			continue
		}

		if r.entries[0] == value {
			return r, nil
		}
	}

	return nil, tracederrors.TracedErrorf("No row found with '%s' as first column value", value)
}

func (s *SpreadSheet) GetMinColumnWithsAsSelectedInOptions(options *SpreadSheetRenderOptions) (columnWidths []int, err error) {
	if options == nil {
		return nil, tracederrors.TracedError("options is nil")
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
		return nil, tracederrors.TracedErrorf("Invalid rowIndex: '%d'", rowIndex)
	}

	rows, err := s.GetRows()
	if err != nil {
		return nil, err
	}

	nRows := len(rows)

	if rowIndex >= nRows {
		return nil, tracederrors.TracedErrorf("Invlaid rowIndex: Index '%d' is invalid for a spreadsheet with '%d' rows.", rowIndex, nRows)
	}

	row = rows[rowIndex]

	if row == nil {
		return nil, tracederrors.TracedError("row is nil after evaluation")
	}

	return row, nil
}

func (s *SpreadSheet) GetRows() (rows []*SpreadSheetRow, err error) {
	if s.rows == nil {
		return nil, tracederrors.TracedErrorf("rows not set")
	}

	if len(s.rows) <= 0 {
		return nil, tracederrors.TracedErrorf("rows has no elements")
	}

	return s.rows, nil
}

func (s *SpreadSheet) GetTitleRow() (TitleRow *SpreadSheetRow, err error) {
	if s.TitleRow == nil {
		return nil, tracederrors.TracedErrorf("TitleRow not set")
	}

	return s.TitleRow, nil
}

func (s *SpreadSheet) MustAddRow(rowEntries []string) {
	err := s.AddRow(rowEntries)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustGetCellValueAsString(rowIndex int, columnIndex int) (cellValue string) {
	cellValue, err := s.GetCellValueAsString(rowIndex, columnIndex)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cellValue
}

func (s *SpreadSheet) MustGetColumnIndexByName(columnName string) (columnIndex int) {
	columnIndex, err := s.GetColumnIndexByName(columnName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return columnIndex
}

func (s *SpreadSheet) MustGetColumnTitleAtIndexAsString(index int) (title string) {
	title, err := s.GetColumnTitleAtIndexAsString(index)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return title
}

func (s *SpreadSheet) MustGetColumnTitlesAsStringSlice() (titles []string) {
	titles, err := s.GetColumnTitlesAsStringSlice()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return titles
}

func (s *SpreadSheet) MustGetMaxColumnWidths() (columnWitdhs []int) {
	columnWitdhs, err := s.GetMaxColumnWidths()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return columnWitdhs
}

func (s *SpreadSheet) MustGetMinColumnWithsAsSelectedInOptions(options *SpreadSheetRenderOptions) (columnWidths []int) {
	columnWidths, err := s.GetMinColumnWithsAsSelectedInOptions(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return columnWidths
}

func (s *SpreadSheet) MustGetNumberOfColumns() (nColumns int) {
	nColumns, err := s.GetNumberOfColumns()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nColumns
}

func (s *SpreadSheet) MustGetNumberOfRows() (nRows int) {
	nRows, err := s.GetNumberOfRows()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nRows
}

func (s *SpreadSheet) MustGetRowByIndex(rowIndex int) (row *SpreadSheetRow) {
	row, err := s.GetRowByIndex(rowIndex)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return row
}

func (s *SpreadSheet) MustGetRows() (rows []*SpreadSheetRow) {
	rows, err := s.GetRows()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return rows
}

func (s *SpreadSheet) MustGetTitleRow() (TitleRow *SpreadSheetRow) {
	TitleRow, err := s.GetTitleRow()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return TitleRow
}

func (s *SpreadSheet) MustPrintAsString(options *SpreadSheetRenderOptions) {
	err := s.PrintAsString(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustRemoveColumnByIndex(columnIndex int) {
	err := s.RemoveColumnByIndex(columnIndex)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustRemoveColumnByName(columnName string) {
	err := s.RemoveColumnByName(columnName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustRenderAsString(options *SpreadSheetRenderOptions) (rendered string) {
	rendered, err := s.RenderAsString(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return rendered
}

func (s *SpreadSheet) MustRenderTitleRowAsString(options *SpreadSheetRenderRowOptions) (rendered string) {
	rendered, err := s.RenderTitleRowAsString(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return rendered
}

func (s *SpreadSheet) MustRenderToStdout(options *SpreadSheetRenderOptions) {
	err := s.RenderToStdout(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustSetColumnTitles(titles []string) {
	err := s.SetColumnTitles(titles)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustSetRows(rows []*SpreadSheetRow) {
	err := s.SetRows(rows)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustSetTitleRow(TitleRow *SpreadSheetRow) {
	err := s.SetTitleRow(TitleRow)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) MustSortByColumnByName(columnName string) {
	err := s.SortByColumnByName(columnName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SpreadSheet) PrintAsString(options *SpreadSheetRenderOptions) (err error) {
	if options == nil {
		return tracederrors.TracedError("options is nil")
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
		return tracederrors.TracedErrorf("columnIndex '%d' is invalid", columnIndex)
	}

	nColumns, err := s.GetNumberOfColumns()
	if err != nil {
		return err
	}

	if columnIndex >= nColumns {
		return tracederrors.TracedErrorf(
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
		return tracederrors.TracedError("columntName is empty string")
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
		return "", tracederrors.TracedErrorNil("options")
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
	renderRowOptions.Prefix = options.Prefix
	renderRowOptions.Suffix = options.Suffix
	renderRowOptions.TitleUnderline = options.TitleUnderline

	rendered = ""
	if !options.SkipTitle {
		toAdd, err := s.RenderTitleRowAsString(renderRowOptions)
		if err != nil {
			return "", err
		}

		rendered += toAdd + "\n"

		if renderRowOptions.TitleUnderline != "" {
			toAdd, err := s.RenderTitleUnderlineAsString(renderRowOptions)
			if err != nil {
				return "", err
			}

			rendered += toAdd + "\n"
		}
	}

	nRows, err := s.GetNumberOfRows()
	if err != nil {
		return "", err
	}

	if nRows > 0 {
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
	}

	return rendered, nil
}

func (s *SpreadSheet) RenderTitleUnderlineAsString(options *SpreadSheetRenderRowOptions) (rendered string, err error) {
	if options == nil {
		return "", tracederrors.TracedError("options is nil")
	}

	titleRow, err := s.GetTitleRow()
	if err != nil {
		return "", err
	}

	if options.Prefix != "" {
		rendered += options.Prefix + " "
	}

	entries, err := titleRow.GetEntries()
	if err != nil {
		return "", err
	}

	for i, e := range entries {
		if i > 0 {
			if options.StringDelimiter == "" {
				rendered += " "
			} else {
				rendered += " " + options.StringDelimiter + " "
			}
		}

		rendered += strings.Repeat("-", len(e))
	}

	if options.Suffix != "" {
		rendered += " " + options.Suffix
	}

	return rendered, nil
}

func (s *SpreadSheet) RenderTitleRowAsString(options *SpreadSheetRenderRowOptions) (rendered string, err error) {
	if options == nil {
		return "", tracederrors.TracedError("options is nil")
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
		return tracederrors.TracedErrorNil("options")
	}

	rendered, err := s.RenderAsString(options)
	if err != nil {
		return err
	}

	rendered = stringsutils.EnsureEndsWithLineBreak(rendered)

	fmt.Print(rendered)

	return nil
}

func (s *SpreadSheet) SetColumnTitles(titles []string) (err error) {
	if len(titles) <= 0 {
		return tracederrors.TracedError("titles is empty slice")
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
		return tracederrors.TracedErrorf("rows is nil")
	}

	if len(rows) <= 0 {
		return tracederrors.TracedErrorf("rows has no elements")
	}

	s.rows = rows

	return nil
}

func (s *SpreadSheet) SetTitleRow(TitleRow *SpreadSheetRow) (err error) {
	if TitleRow == nil {
		return tracederrors.TracedErrorf("TitleRow is nil")
	}

	s.TitleRow = TitleRow

	return nil
}

func (s *SpreadSheet) SortByColumnByName(columnName string) (err error) {
	if columnName == "" {
		return tracederrors.TracedError("columnName is empty string")
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
			logging.LogGoErrorFatal(err)
		}

		jRow, err := s.GetRowByIndex(j)
		if err != nil {
			logging.LogGoErrorFatal(err)
		}

		iValue, err := iRow.GetColumnValueAsString(columnIndex)
		if err != nil {
			logging.LogGoErrorFatal(err)
		}

		jValue, err := jRow.GetColumnValueAsString(columnIndex)
		if err != nil {
			logging.LogGoErrorFatal(err)
		}

		return iValue < jValue
	})

	err = s.SetRows(rows)
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheet) UpdateRowFoundByFirstColumnValue(searchValue string, cellIndex int, newValue string) (err error) {
	if cellIndex < 0 {
		return tracederrors.TracedErrorf("Invalid cellIndex '%d' to update", cellIndex)
	}

	row, err := s.GetRowByFirstColumnValue(searchValue)
	if err != nil {
		return err
	}

	entries, err := row.GetEntries()
	if err != nil {
		return err
	}

	nEntries := len(entries)
	if nEntries <= cellIndex {
		return tracederrors.TracedErrorf("cellIndex to update '%d' is to high. There are only '%d' elements in the row.", cellIndex, nEntries)
	}

	entries[cellIndex] = newValue

	return nil
}

func (s *SpreadSheet) MustGetRowByIndexAsStringSlice(index int) (values []string) {
	values, err := s.GetRowByIndexAsStringSlice(index)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return values
}

func (s *SpreadSheet) GetRowByIndexAsStringSlice(index int) (values []string, err error) {
	row, err := s.GetRowByIndex(index)
	if err != nil {
		return nil, err
	}

	return row.GetEntries()
}

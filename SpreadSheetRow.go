package asciichgolangpublic

import (
	"strings"
)

type SpreadSheetRow struct {
	entries []string
}

func NewSpreadSheetRow() (s *SpreadSheetRow) {
	return new(SpreadSheetRow)
}

func (s *SpreadSheetRow) GetColumnValueAsString(columnIndex int) (columnValue string, err error) {
	if columnIndex < 0 {
		return "", TracedErrorf("Invalid columnIndex: '%d'", columnIndex)
	}

	entries, err := s.GetEntries()
	if err != nil {
		return "", err
	}

	nEntries := len(entries)
	if columnIndex >= nEntries {
		return "", TracedErrorf(
			"Invalid columnIndex '%d' for a spread sheet with '%d' columns",
			columnIndex,
			nEntries,
		)
	}

	columnValue = entries[columnIndex]

	return columnValue, nil
}

func (s *SpreadSheetRow) GetColumnWidths() (columnWidths []int, err error) {
	nColumns, err := s.GetNumberOfEntries()
	if err != nil {
		return nil, err
	}

	columnWidths = []int{}
	for i := 0; i < nColumns; i++ {
		content, err := s.GetColumnValueAsString(i)
		if err != nil {
			return nil, err
		}

		columnWidths = append(columnWidths, len(content))
	}

	return columnWidths, nil
}

func (s *SpreadSheetRow) GetEntries() (entries []string, err error) {
	if s.entries == nil {
		return nil, TracedErrorf("entries not set")
	}

	if len(s.entries) <= 0 {
		return nil, TracedErrorf("entries has no elements")
	}

	return s.entries, nil
}

func (s *SpreadSheetRow) GetNumberOfEntries() (nEntries int, err error) {
	entries, err := s.GetEntries()
	if err != nil {
		return -1, err
	}

	return len(entries), nil
}

func (s *SpreadSheetRow) MustGetColumnValueAsString(columnIndex int) (columnValue string) {
	columnValue, err := s.GetColumnValueAsString(columnIndex)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return columnValue
}

func (s *SpreadSheetRow) MustGetColumnWidths() (columnWidths []int) {
	columnWidths, err := s.GetColumnWidths()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return columnWidths
}

func (s *SpreadSheetRow) MustGetEntries() (entries []string) {
	entries, err := s.GetEntries()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return entries
}

func (s *SpreadSheetRow) MustGetNumberOfEntries() (nEntries int) {
	nEntries, err := s.GetNumberOfEntries()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nEntries
}

func (s *SpreadSheetRow) MustRemoveElementAtIndex(index int) {
	err := s.RemoveElementAtIndex(index)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheetRow) MustRenderAsString(options *SpreadSheetRenderRowOptions) (rendered string) {
	rendered, err := s.RenderAsString(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return rendered
}

func (s *SpreadSheetRow) MustSetEntries(entries []string) {
	err := s.SetEntries(entries)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheetRow) RemoveElementAtIndex(index int) (err error) {
	if index < 0 {
		return TracedErrorf("Index '%d' is invalid.", index)
	}

	entries, err := s.GetEntries()
	if err != nil {
		return err
	}

	entries = Slices().RemoveStringEntryAtIndex(entries, index)

	err = s.SetEntries(entries)
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetRow) RenderAsString(options *SpreadSheetRenderRowOptions) (rendered string, err error) {
	if options == nil {
		return "", TracedError("options is nil")
	}

	entries, err := s.GetEntries()
	if err != nil {
		return "", err
	}

	var delimiter string = ""
	if options.IsStringDelimiterSet() {
		delimiter, err = options.GetStringDelimiter()
		if err != nil {
			return "", err
		}

		delimiter = " " + delimiter + " "
	} else {
		delimiter = " "
	}

	if options.IsMinColumnWidthsSet() {
		minColumnWidth, err := options.GetMinColumnWidths()
		if err != nil {
			return "", err
		}

		nEntries := len(entries)
		nCloumnWidths := len(minColumnWidth)
		if nEntries != nCloumnWidths {
			return "", TracedErrorf("nEntries = '%d' != nCloumnWidths = '%d'", nEntries, nCloumnWidths)
		}

		entriesFilled := []string{}
		for i, entry := range entries {
			entriesFilled = append(
				entriesFilled,
				Strings().RightFillWithSpaces(entry, minColumnWidth[i]),
			)
		}

		entries = entriesFilled
	}

	rendered = strings.Join(entries, delimiter)
	return rendered, nil
}

func (s *SpreadSheetRow) SetEntries(entries []string) (err error) {
	if entries == nil {
		return TracedErrorf("entries is nil")
	}

	if len(entries) <= 0 {
		return TracedErrorf("entries has no elements")
	}

	s.entries = entries

	return nil
}

package spreadsheet

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type SpreadSheetRenderRowOptions struct {
	MinColumnWidths []int
	StringDelimiter string
	Verbose         bool
	Prefix          string
	Suffix          string
	TitleUnderline  string
}

func NewSpreadSheetRenderRowOptions() (s *SpreadSheetRenderRowOptions) {
	return new(SpreadSheetRenderRowOptions)
}

func (s *SpreadSheetRenderRowOptions) GetMinColumnWidths() (minColumnWidths []int, err error) {
	if s.MinColumnWidths == nil {
		return nil, tracederrors.TracedErrorf("MinColumnWidths not set")
	}

	if len(s.MinColumnWidths) <= 0 {
		return nil, tracederrors.TracedErrorf("MinColumnWidths has no elements")
	}

	return s.MinColumnWidths, nil
}

func (s *SpreadSheetRenderRowOptions) GetStringDelimiter() (stringDelimiter string, err error) {
	if s.StringDelimiter == "" {
		return "", tracederrors.TracedErrorf("StringDelimiter not set")
	}

	return s.StringDelimiter, nil
}

func (s *SpreadSheetRenderRowOptions) GetVerbose() (verbose bool, err error) {

	return s.Verbose, nil
}

func (s *SpreadSheetRenderRowOptions) IsMinColumnWidthsSet() (isSet bool) {
	return len(s.MinColumnWidths) > 0
}

func (s *SpreadSheetRenderRowOptions) IsStringDelimiterSet() (isSet bool) {
	return s.StringDelimiter != ""
}

func (s *SpreadSheetRenderRowOptions) SetMinColumnWidths(minColumnWidths []int) (err error) {
	if minColumnWidths == nil {
		return tracederrors.TracedErrorf("minColumnWidths is nil")
	}

	if len(minColumnWidths) <= 0 {
		return tracederrors.TracedErrorf("minColumnWidths has no elements")
	}

	s.MinColumnWidths = minColumnWidths

	return nil
}

func (s *SpreadSheetRenderRowOptions) SetStringDelimiter(stringDelimiter string) (err error) {
	if stringDelimiter == "" {
		return tracederrors.TracedErrorf("stringDelimiter is empty string")
	}

	s.StringDelimiter = stringDelimiter

	return nil
}

func (s *SpreadSheetRenderRowOptions) SetVerbose(verbose bool) (err error) {
	s.Verbose = verbose

	return nil
}

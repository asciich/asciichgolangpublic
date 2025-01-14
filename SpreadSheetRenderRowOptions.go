package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type SpreadSheetRenderRowOptions struct {
	MinColumnWidths []int
	StringDelimiter string
	Verbose         bool
}

func NewSpreadSheetRenderRowOptions() (s *SpreadSheetRenderRowOptions) {
	return new(SpreadSheetRenderRowOptions)
}

func (s *SpreadSheetRenderRowOptions) GetMinColumnWidths() (minColumnWidths []int, err error) {
	if s.MinColumnWidths == nil {
		return nil, errors.TracedErrorf("MinColumnWidths not set")
	}

	if len(s.MinColumnWidths) <= 0 {
		return nil, errors.TracedErrorf("MinColumnWidths has no elements")
	}

	return s.MinColumnWidths, nil
}

func (s *SpreadSheetRenderRowOptions) GetStringDelimiter() (stringDelimiter string, err error) {
	if s.StringDelimiter == "" {
		return "", errors.TracedErrorf("StringDelimiter not set")
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

func (s *SpreadSheetRenderRowOptions) MustGetMinColumnWidths() (minColumnWidths []int) {
	minColumnWidths, err := s.GetMinColumnWidths()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return minColumnWidths
}

func (s *SpreadSheetRenderRowOptions) MustGetStringDelimiter() (stringDelimiter string) {
	stringDelimiter, err := s.GetStringDelimiter()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return stringDelimiter
}

func (s *SpreadSheetRenderRowOptions) MustGetVerbose() (verbose bool) {
	verbose, err := s.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (s *SpreadSheetRenderRowOptions) MustSetMinColumnWidths(minColumnWidths []int) {
	err := s.SetMinColumnWidths(minColumnWidths)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SpreadSheetRenderRowOptions) MustSetStringDelimiter(stringDelimiter string) {
	err := s.SetStringDelimiter(stringDelimiter)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SpreadSheetRenderRowOptions) MustSetVerbose(verbose bool) {
	err := s.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SpreadSheetRenderRowOptions) SetMinColumnWidths(minColumnWidths []int) (err error) {
	if minColumnWidths == nil {
		return errors.TracedErrorf("minColumnWidths is nil")
	}

	if len(minColumnWidths) <= 0 {
		return errors.TracedErrorf("minColumnWidths has no elements")
	}

	s.MinColumnWidths = minColumnWidths

	return nil
}

func (s *SpreadSheetRenderRowOptions) SetStringDelimiter(stringDelimiter string) (err error) {
	if stringDelimiter == "" {
		return errors.TracedErrorf("stringDelimiter is empty string")
	}

	s.StringDelimiter = stringDelimiter

	return nil
}

func (s *SpreadSheetRenderRowOptions) SetVerbose(verbose bool) (err error) {
	s.Verbose = verbose

	return nil
}

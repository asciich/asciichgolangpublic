package spreadsheet

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type SpreadSheetRenderOptions struct {
	SkipTitle                 bool
	StringDelimiter           string
	Verbose                   bool
	SameColumnWidthForAllRows bool
	Prefix                    string
	Suffix                    string
	TitleUnderline            string
}

func NewSpreadSheetRenderOptions() (s *SpreadSheetRenderOptions) {
	return new(SpreadSheetRenderOptions)
}

func (s *SpreadSheetRenderOptions) GetSameColumnWidthForAllRows() (sameColumnWidthForAllRows bool, err error) {

	return s.SameColumnWidthForAllRows, nil
}

func (s *SpreadSheetRenderOptions) GetSkipTitle() (skipTitle bool, err error) {

	return s.SkipTitle, nil
}

func (s *SpreadSheetRenderOptions) GetStringDelimiter() (stringDelimiter string, err error) {
	if s.StringDelimiter == "" {
		return "", tracederrors.TracedErrorf("StringDelimiter not set")
	}

	return s.StringDelimiter, nil
}

func (s *SpreadSheetRenderOptions) GetVerbose() (verbose bool, err error) {

	return s.Verbose, nil
}

func (s *SpreadSheetRenderOptions) SetSameColumnWidthForAllRows(sameColumnWidthForAllRows bool) (err error) {
	s.SameColumnWidthForAllRows = sameColumnWidthForAllRows

	return nil
}

func (s *SpreadSheetRenderOptions) SetSkipTitle(skipTitle bool) (err error) {
	s.SkipTitle = skipTitle

	return nil
}

func (s *SpreadSheetRenderOptions) SetStringDelimiter(stringDelimiter string) (err error) {
	if stringDelimiter == "" {
		return tracederrors.TracedErrorf("stringDelimiter is empty string")
	}

	s.StringDelimiter = stringDelimiter

	return nil
}

func (s *SpreadSheetRenderOptions) SetVerbose(verbose bool) (err error) {
	s.Verbose = verbose

	return nil
}

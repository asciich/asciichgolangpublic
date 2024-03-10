package asciichgolangpublic

type SpreadSheetRenderOptions struct {
	SkipTitle                 bool
	StringDelimiter           string
	Verbose                   bool
	SameColumnWidthForAllRows bool
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
		return "", TracedErrorf("StringDelimiter not set")
	}

	return s.StringDelimiter, nil
}

func (s *SpreadSheetRenderOptions) GetVerbose() (verbose bool, err error) {

	return s.Verbose, nil
}

func (s *SpreadSheetRenderOptions) MustGetSameColumnWidthForAllRows() (sameColumnWidthForAllRows bool) {
	sameColumnWidthForAllRows, err := s.GetSameColumnWidthForAllRows()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sameColumnWidthForAllRows
}

func (s *SpreadSheetRenderOptions) MustGetSkipTitle() (skipTitle bool) {
	skipTitle, err := s.GetSkipTitle()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return skipTitle
}

func (s *SpreadSheetRenderOptions) MustGetStringDelimiter() (stringDelimiter string) {
	stringDelimiter, err := s.GetStringDelimiter()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stringDelimiter
}

func (s *SpreadSheetRenderOptions) MustGetVerbose() (verbose bool) {
	verbose, err := s.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (s *SpreadSheetRenderOptions) MustSetSameColumnWidthForAllRows(sameColumnWidthForAllRows bool) {
	err := s.SetSameColumnWidthForAllRows(sameColumnWidthForAllRows)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheetRenderOptions) MustSetSkipTitle(skipTitle bool) {
	err := s.SetSkipTitle(skipTitle)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheetRenderOptions) MustSetStringDelimiter(stringDelimiter string) {
	err := s.SetStringDelimiter(stringDelimiter)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SpreadSheetRenderOptions) MustSetVerbose(verbose bool) {
	err := s.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
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
		return TracedErrorf("stringDelimiter is empty string")
	}

	s.StringDelimiter = stringDelimiter

	return nil
}

func (s *SpreadSheetRenderOptions) SetVerbose(verbose bool) (err error) {
	s.Verbose = verbose

	return nil
}

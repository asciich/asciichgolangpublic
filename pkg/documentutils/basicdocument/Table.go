package basicdocument

import (
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/spreadsheet"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Table struct {
	spreadsheet *spreadsheet.SpreadSheet
}

func NewTable() (t *Table) {
	t = new(Table)
	t.spreadsheet = spreadsheet.NewSpreadSheet()
	return t
}

func GetNewTable() (t *Table, err error) {
	t = NewTable()

	return t, nil
}

func (t *Table) GetSpreadSheet() (spreadSheet *spreadsheet.SpreadSheet, err error) {
	if t.spreadsheet == nil {
		return nil, tracederrors.TracedError("SpreadSheet not set")
	}

	return t.spreadsheet, nil
}

func (t *Table) GetPlainText() (plain string, err error) {
	return "table", nil
}

func (t *Table) MustRenderAsMarkdownString() (rendered string) {
	rendered, err := t.RenderAsMarkDownString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return rendered
}

func (t *Table) RenderAsMarkDownString() (rendered string, err error) {
	sheet, err := t.GetSpreadSheet()
	if err != nil {
		return "", err
	}

	rendered, err = sheet.RenderAsString(
		&spreadsheet.SpreadSheetRenderOptions{
			StringDelimiter:           "|",
			SameColumnWidthForAllRows: true,
			Prefix:                    "|",
			Suffix:                    "|",
			TitleUnderline:            "-",
		},
	)
	if err != nil {
		return "", err
	}

	return rendered, nil
}

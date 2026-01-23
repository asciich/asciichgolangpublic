package basicdocument

import (
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/documentinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/spreadsheet"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type BasicDocument struct {
	elements []documentinterfaces.Element
}

func NewBasicDocument() (d *BasicDocument) {
	return new(BasicDocument)
}

func (d BasicDocument) GetElements() (elementes []documentinterfaces.Element) {
	return d.elements
}

func (d *BasicDocument) AddElement(element documentinterfaces.Element) error {
	if element == nil {
		return tracederrors.TracedErrorNil("element")
	}

	d.elements = append(d.elements, element)

	return nil
}

func (d *BasicDocument) AddTable() (spreadsheet *spreadsheet.SpreadSheet, err error) {
	toAdd, err := GetNewTable()
	if err != nil {
		return nil, err
	}

	err = d.AddElement(toAdd)
	if err != nil {
		return nil, err
	}

	return toAdd.GetSpreadSheet()
}

func (d *BasicDocument) AddTitleByString(title string) error {
	if title == "" {
		return tracederrors.TracedErrorEmptyString("title")
	}

	toAdd, err := GetNewTitleByString(title)
	if err != nil {
		return err
	}

	return d.AddElement(toAdd)
}

func (d *BasicDocument) AddSubTitleByString(subtitle string) error {
	if subtitle == "" {
		return tracederrors.TracedErrorEmptyString("subtitle")
	}

	toAdd, err := GetNewSubTitleByString(subtitle)
	if err != nil {
		return err
	}

	err = d.AddElement(toAdd)
	if err != nil {
		return err
	}

	return nil
}

func (d *BasicDocument) AddSubSubTitleByString(subsubtitle string) error {
	if subsubtitle == "" {
		return tracederrors.TracedErrorEmptyString("subsubtitle")
	}

	toAdd, err := GetNewSubSubTitleByString(subsubtitle)
	if err != nil {
		return err
	}

	return d.AddElement(toAdd)
}

func (d *BasicDocument) AddSubSubSubTitleByString(subsubsubtitle string) error {
	if subsubsubtitle == "" {
		return tracederrors.TracedErrorEmptyString("subsubsubtitle")
	}

	toAdd, err := GetNewSubSubSubTitleByString(subsubsubtitle)
	if err != nil {
		return err
	}

	return d.AddElement(toAdd)
}

func (d *BasicDocument) AddTextByString(text string) error {
	if text == "" {
		return tracederrors.TracedErrorEmptyString("text")
	}

	toAdd, err := GetNewTextByString(text)
	if err != nil {
		return err
	}

	return d.AddElement(toAdd)
}

func (d *BasicDocument) AddVerbatimByString(verbatim string) error {
	toAdd, err := GetNewVerbatimByString(verbatim)
	if err != nil {
		return err
	}

	return d.AddElement(toAdd)
}

// Add a code block in the given language.
// Use an empty string for the language if not specified.
func (d *BasicDocument) AddCodeBlockByString(code string, language string) error {
	toAdd, err := GetNewCodeBlockByString(code, language)
	if err != nil {
		return err
	}

	return d.AddElement(toAdd)
}

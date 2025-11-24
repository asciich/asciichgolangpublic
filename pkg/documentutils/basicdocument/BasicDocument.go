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

func (d *BasicDocument) AddElement(element documentinterfaces.Element) (err error) {
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

func (d *BasicDocument) AddTitleByString(title string) (err error) {
	if title == "" {
		return tracederrors.TracedErrorEmptyString("title")
	}

	toAdd, err := GetNewTitleByString(title)
	if err != nil {
		return err
	}

	err = d.AddElement(toAdd)
	if err != nil {
		return err
	}

	return nil
}

func (d *BasicDocument) AddSubTitleByString(subtitle string) (err error) {
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

func (d *BasicDocument) AddSubSubTitleByString(subsubtitle string) (err error) {
	if subsubtitle == "" {
		return tracederrors.TracedErrorEmptyString("subsubtitle")
	}

	toAdd, err := GetNewSubSubTitleByString(subsubtitle)
	if err != nil {
		return err
	}

	err = d.AddElement(toAdd)
	if err != nil {
		return err
	}

	return nil
}

func (d *BasicDocument) AddSubSubSubTitleByString(subsubsubtitle string) (err error) {
	if subsubsubtitle == "" {
		return tracederrors.TracedErrorEmptyString("subsubsubtitle")
	}

	toAdd, err := GetNewSubSubSubTitleByString(subsubsubtitle)
	if err != nil {
		return err
	}

	err = d.AddElement(toAdd)
	if err != nil {
		return err
	}

	return nil
}

func (d *BasicDocument) AddTextByString(text string) (err error) {
	if text == "" {
		return tracederrors.TracedErrorEmptyString("text")
	}

	toAdd, err := GetNewTextByString(text)
	if err != nil {
		return err
	}

	err = d.AddElement(toAdd)
	if err != nil {
		return err
	}

	return nil
}

func (d *BasicDocument) AddVerbatimByString(verbatim string) (err error) {
	toAdd, err := GetNewVerbatimByString(verbatim)
	if err != nil {
		return err
	}

	err = d.AddElement(toAdd)
	if err != nil {
		return err
	}

	return nil
}

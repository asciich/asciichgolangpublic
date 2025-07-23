package documentbase

import (
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/spreadsheet"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

type DocumentBase struct {
	elements []Element
}

func NewDocumentBase() (d *DocumentBase) {
	return new(DocumentBase)
}

func (d DocumentBase) GetElements() (elementes []Element) {
	return d.elements
}

func (d *DocumentBase) AddElement(element Element) (err error) {
	if element == nil {
		return tracederrors.TracedErrorNil("element")
	}

	d.elements = append(d.elements, element)

	return nil
}

func (d *DocumentBase) AddTable() (spreadsheet *spreadsheet.SpreadSheet, err error) {
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

func (d *DocumentBase) AddTitleByString(title string) (err error) {
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

func (d *DocumentBase) AddSubTitleByString(subtitle string) (err error) {
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

func (d *DocumentBase) AddSubSubTitleByString(subsubtitle string) (err error) {
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

func (d *DocumentBase) AddSubSubSubTitleByString(subsubsubtitle string) (err error) {
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

func (d *DocumentBase) AddTextByString(text string) (err error) {
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

func (d *DocumentBase) RenderAsString() (rendered string, err error) {
	return "", tracederrors.TracedError("DocumentBase does not implement any renderer. You have to implement your own renderer.")
}

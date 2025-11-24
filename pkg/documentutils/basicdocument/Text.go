package basicdocument

type Text struct {
	ElementBase
}

func NewText() (t *Text) {
	return new(Text)
}

func GetNewTextByString(text string) (t *Text, err error) {
	t = NewText()

	err = t.SetPlainText(text)
	if err != nil {
		return nil, err
	}

	return t, nil
}

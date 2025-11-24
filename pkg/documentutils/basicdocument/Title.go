package basicdocument

type Title struct {
	ElementBase
}

func NewTitle() (t *Title) {
	return new(Title)
}

func GetNewTitleByString(title string) (t *Title, err error) {
	t = NewTitle()

	err = t.SetPlainText(title)
	if err != nil {
		return nil, err
	}

	return t, nil
}

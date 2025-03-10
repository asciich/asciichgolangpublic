package documentbase

type SubSubTitle struct {
	ElementBase
}

func NewSubSubTitle() (t *SubSubTitle) {
	return new(SubSubTitle)
}

func GetNewSubSubTitleByString(subtitle string) (t *SubSubTitle, err error) {
	t = NewSubSubTitle()

	err = t.SetPlainText(subtitle)
	if err != nil {
		return nil, err
	}

	return t, nil
}

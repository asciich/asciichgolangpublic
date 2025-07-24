package documentbase

type SubSubSubTitle struct {
	ElementBase
}

func NewSubSubSubTitle() (t *SubSubSubTitle) {
	return new(SubSubSubTitle)
}

func GetNewSubSubSubTitleByString(subtitle string) (t *SubSubSubTitle, err error) {
	t = NewSubSubSubTitle()

	err = t.SetPlainText(subtitle)
	if err != nil {
		return nil, err
	}

	return t, nil
}

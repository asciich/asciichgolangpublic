package documentbase

type SubTitle struct {
	ElementBase
}

func NewSubTitle() (t *SubTitle) {
	return new(SubTitle)
}

func GetNewSubTitleByString(subtitle string) (t *SubTitle, err error) {
	t = NewSubTitle()

	err = t.SetPlainText(subtitle)
	if err != nil {
		return nil, err
	}

	return t, nil
}

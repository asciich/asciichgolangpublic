// Verbatim.go

package basicdocument

type Verbatim struct {
	ElementBase
}

func NewVerbatim() (v *Verbatim) {
	return new(Verbatim)
}

func GetNewVerbatimByString(content string) (*Verbatim, error) {
	v := NewVerbatim()

	if content != "" {
		err := v.SetPlainText(content)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

package documentbase

import (
	"strings"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

type ElementBase struct {
	plainText string
}

func NewElementBase() (e *ElementBase) {
	return new(ElementBase)
}

func (e *ElementBase) SetPlainText(plainText string) (err error) {
	plainText = strings.TrimSpace(plainText)

	if plainText == "" {
		return tracederrors.TracedErrorEmptyString("plainText")
	}

	e.plainText = plainText

	return nil
}

func (e ElementBase) GetPlainText() (plainText string, err error) {
	if e.plainText == "" {
		return "", tracederrors.TracedError("plainText not set")
	}

	return e.plainText, nil
}

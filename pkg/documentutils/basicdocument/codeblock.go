package basicdocument

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CodeBlock struct {
	ElementBase
	langauge string
}

func NewCodeBlock() (v *CodeBlock) {
	return new(CodeBlock)
}

func GetNewCodeBlockByString(content string, language string) (*CodeBlock, error) {
	c := NewCodeBlock()

	if content != "" {
		err := c.SetPlainText(content)
		if err != nil {
			return nil, err
		}
	}

	if language != "" {
		err := c.SetLanguage(language)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *CodeBlock) SetLanguage(language string) error {
	if language == "" {
		return tracederrors.TracedErrorEmptyString("language")
	}

	c.langauge = language

	return nil
}

func (c *CodeBlock) GetLanguage() (string, error) {
	if c.langauge == "" {
		return "", tracederrors.TracedError("Langauge not set")
	}

	return c.langauge, nil
}

func (c *CodeBlock) GetLanguageOrEmptyIfUnset() string {
	return c.langauge
}

func (c *CodeBlock) IsLanguageSet() bool {
	return c.langauge != ""
}

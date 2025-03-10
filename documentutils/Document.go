package documentutils

import "github.com/asciich/asciichgolangpublic/documentutils/documentbase"

type Document interface {
	AddTitleByString(title string) (err error)
	GetElements() (elements []documentbase.Element)
	RenderAsString() (rendered string, err error)
}

package markdowndocument

import (
	"reflect"

	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/documentutils/documentbase"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type MarkDownDocument struct {
	documentbase.DocumentBase
}

func NewMarkDownDocument() (m *MarkDownDocument) {
	return new(MarkDownDocument)
}

func (m MarkDownDocument) RenderAsString() (rendered string, err error) {
	extendRendered := func(r string) string {
		if r == "" {
			return r
		}

		return r + "\n"
	}

	for _, e := range m.GetElements() {
		plainTest, err := e.GetPlainText()
		if err != nil {
			return "", err
		}

		switch e := e.(type) {
		case *documentbase.Text:
			rendered = extendRendered(rendered) + stringsutils.EnsureEndsWithExactlyOneLineBreak(plainTest)
		case *documentbase.Title:
			rendered = extendRendered(rendered) + "# " + stringsutils.EnsureEndsWithExactlyOneLineBreak(plainTest)
		case *documentbase.SubTitle:
			rendered = extendRendered(rendered) + "## " + stringsutils.EnsureEndsWithExactlyOneLineBreak(plainTest)
		case *documentbase.SubSubTitle:
			rendered = extendRendered(rendered) + "### " + stringsutils.EnsureEndsWithExactlyOneLineBreak(plainTest)
		case *documentbase.SubSubSubTitle:
			rendered = extendRendered(rendered) + "#### " + stringsutils.EnsureEndsWithExactlyOneLineBreak(plainTest)
		case *documentbase.Table:
			table, err := e.RenderAsMarkDownString()
			if err != nil {
				return "", err
			}
			rendered = extendRendered(rendered) + stringsutils.EnsureEndsWithExactlyOneLineBreak(table)
		default:
			return "", tracederrors.TracedErrorf("Unknown element type to render: %s", reflect.TypeOf(e))
		}
	}

	rendered = stringsutils.EnsureEndsWithExactlyOneLineBreak(rendered)

	return rendered, nil
}

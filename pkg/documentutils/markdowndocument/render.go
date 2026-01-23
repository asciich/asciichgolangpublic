package markdowndocument

import (
	"reflect"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/basicdocument"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/documentinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func RenderAsString(document documentinterfaces.Document) (rendered string, err error) {
	extendRendered := func(r string) string {
		if r == "" {
			return r
		}

		return r + "\n"
	}

	for _, e := range document.GetElements() {
		plainTest, err := e.GetPlainText()
		if err != nil {
			return "", err
		}

		switch e := e.(type) {
		case *basicdocument.Text:
			rendered = extendRendered(rendered) + stringsutils.EnsureEndsWithExactlyOneLineBreak(plainTest)
		case *basicdocument.Title:
			rendered = extendRendered(rendered) + "# " + stringsutils.EnsureEndsWithExactlyOneLineBreak(plainTest)
		case *basicdocument.SubTitle:
			rendered = extendRendered(rendered) + "## " + stringsutils.EnsureEndsWithExactlyOneLineBreak(plainTest)
		case *basicdocument.SubSubTitle:
			rendered = extendRendered(rendered) + "### " + stringsutils.EnsureEndsWithExactlyOneLineBreak(plainTest)
		case *basicdocument.SubSubSubTitle:
			rendered = extendRendered(rendered) + "#### " + stringsutils.EnsureEndsWithExactlyOneLineBreak(plainTest)
		case *basicdocument.Table:
			table, err := e.RenderAsMarkDownString()
			if err != nil {
				return "", err
			}
			rendered = extendRendered(rendered) + stringsutils.EnsureEndsWithExactlyOneLineBreak(table)
		case *basicdocument.Verbatim:
			rendered = extendRendered(rendered) + "```\n" + plainTest + "\n```\n"
		case *basicdocument.CodeBlock:
			language := e.GetLanguageOrEmptyIfUnset()
			rendered = extendRendered(rendered) + "```" + language + "\n" + plainTest + "\n```\n"
		default:
			return "", tracederrors.TracedErrorf("Unknown element type to render: %s", reflect.TypeOf(e))
		}
	}

	rendered = stringsutils.EnsureEndsWithExactlyOneLineBreak(rendered)

	return rendered, nil
}

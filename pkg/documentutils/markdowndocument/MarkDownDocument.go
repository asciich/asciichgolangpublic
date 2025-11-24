package markdowndocument

import (
	"reflect"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/documentinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/basicdocument"
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
		default:
			return "", tracederrors.TracedErrorf("Unknown element type to render: %s", reflect.TypeOf(e))
		}
	}

	rendered = stringsutils.EnsureEndsWithExactlyOneLineBreak(rendered)

	return rendered, nil
}

// ParseFromString parses a markdown string and populates the document
func ParseFromString(markdown string) (documentinterfaces.Document, error) {
	lines := strings.SplitSeq(markdown, "\n")
	var inVerbatim bool
	var verbatimContent []string

	ret := basicdocument.NewBasicDocument()

	for line := range lines {
		if line == "```" {
			if inVerbatim {
				// End of verbatim block
				if err := ret.AddVerbatimByString(strings.Join(verbatimContent, "\n")); err != nil {
					return nil, err
				}
				verbatimContent = nil
				inVerbatim = false
			} else {
				// Start of verbatim block
				inVerbatim = true
			}
			continue
		}

		if inVerbatim {
			verbatimContent = append(verbatimContent, line)
			continue
		}

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#### ") {
			if err := ret.AddSubSubSubTitleByString(strings.TrimPrefix(line, "#### ")); err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(line, "### ") {
			if err := ret.AddSubSubTitleByString(strings.TrimPrefix(line, "### ")); err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(line, "## ") {
			if err := ret.AddSubTitleByString(strings.TrimPrefix(line, "## ")); err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(line, "# ") {
			if err := ret.AddTitleByString(strings.TrimPrefix(line, "# ")); err != nil {
				return nil, err
			}
		} else if !strings.HasPrefix(line, "|") {
			if err := ret.AddTextByString(line); err != nil {
				return nil, err
			}
		}
	}

	return ret, nil
}

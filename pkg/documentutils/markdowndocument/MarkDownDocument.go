package markdowndocument

import (
	"reflect"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/documentbase"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
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
		case *documentbase.Verbatim:
			rendered = extendRendered(rendered) + "```\n" + plainTest + "\n```\n"
		default:
			return "", tracederrors.TracedErrorf("Unknown element type to render: %s", reflect.TypeOf(e))
		}
	}

	rendered = stringsutils.EnsureEndsWithExactlyOneLineBreak(rendered)

	return rendered, nil
}

// ParseFromString parses a markdown string and populates the document
func (d *MarkDownDocument) ParseFromString(markdown string) error {
	lines := strings.SplitSeq(markdown, "\n")
	var inVerbatim bool
	var verbatimContent []string

	for line := range lines {
		if line == "```" {
			if inVerbatim {
				// End of verbatim block
				if err := d.AddVerbatimByString(strings.Join(verbatimContent, "\n")); err != nil {
					return err
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
			if err := d.AddSubSubSubTitleByString(strings.TrimPrefix(line, "#### ")); err != nil {
				return err
			}
		} else if strings.HasPrefix(line, "### ") {
			if err := d.AddSubSubTitleByString(strings.TrimPrefix(line, "### ")); err != nil {
				return err
			}
		} else if strings.HasPrefix(line, "## ") {
			if err := d.AddSubTitleByString(strings.TrimPrefix(line, "## ")); err != nil {
				return err
			}
		} else if strings.HasPrefix(line, "# ") {
			if err := d.AddTitleByString(strings.TrimPrefix(line, "# ")); err != nil {
				return err
			}
		} else if !strings.HasPrefix(line, "|") {
			if err := d.AddTextByString(line); err != nil {
				return err
			}
		}
	}

	return nil
}

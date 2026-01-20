package htmldocument

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/documentutils/basicdocument"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/documentinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/htmlutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"golang.org/x/net/html"
)

func ParseString(htmlString string) (documentinterfaces.Document, error) {
	htmlString = strings.TrimSpace(htmlString)
	if htmlString == "" {
		return nil, tracederrors.TracedErrorEmptyString("htmlString")
	}

	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse htmlString as HTML: %w", err)
	}

	body, err := htmlutils.FindBodyNode(doc)
	if err != nil {
		if htmlutils.IsErrNoHtmlBodyFound(err) {
			return basicdocument.NewBasicDocument(), nil
		}
	}

	document := basicdocument.NewBasicDocument()
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		// Filter for ElementNodes to avoid printing empty whitespace/newlines
		if c.Type == html.ElementNode {
			tag := strings.ToLower(c.Data)
			text := htmlutils.GetNodeText(c)

			switch tag {
			case "h1":
				document.AddTitleByString(text)
			case "h2":
				document.AddSubTitleByString(text)
			case "p":
				document.AddTextByString(text)
			default:
				return nil, tracederrors.TracedErrorf("Unknown tag '%s'", tag)
			}
		}
	}

	return document, nil
}

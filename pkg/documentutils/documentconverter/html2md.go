package documentconverter

import (
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/htmldocument"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/markdowndocument"
)

// Converts a HTML document string into a markdown string.
func HtmlStringToMdString(html string) (string, error) {
	document, err := htmldocument.ParseString(html)
	if err != nil {
		return "", err
	}

	return markdowndocument.RenderAsString(document)
}

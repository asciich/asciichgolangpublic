package htmldocument

import "github.com/asciich/asciichgolangpublic/pkg/documentutils/documentbase"

type HtmlDocument struct {
	documentbase.DocumentBase
}

func NewHtmlDocument() (h *HtmlDocument) {
	return new(HtmlDocument)
}

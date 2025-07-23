package htmldocument

import "gitlab.asciich.ch/tools/asciichgolangpublic.git/documentutils/documentbase"

type HtmlDocument struct {
	documentbase.DocumentBase
}

func NewHtmlDocument() (h *HtmlDocument) {
	return new(HtmlDocument)
}

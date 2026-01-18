package htmldocument

import (
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/documentutils/documentinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func RenderAsString(document documentinterfaces.Document) (string, error) {
	if document == nil {
		return "", tracederrors.TracedErrorEmptyString("document")
	}

	rendered := new(strings.Builder)
	fmt.Fprintf(rendered, "<html>\n")
	fmt.Fprintf(rendered, "<body>\n")
	fmt.Fprintf(rendered, "</body>\n")
	fmt.Fprintf(rendered, "</html>\n")

	return rendered.String(), nil
}

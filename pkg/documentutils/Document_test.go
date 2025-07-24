package documentutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/documentbase"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/htmldocument"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/markdowndocument"
)

func Test_ImplementDocumentInterface(t *testing.T) {
	t.Run("DocumentBase", func(t *testing.T) {
		var d Document = documentbase.NewDocumentBase()

		require.Len(t, d.GetElements(), 0)
	})

	t.Run("HTMLDocument", func(t *testing.T) {
		var d Document = htmldocument.NewHtmlDocument()

		require.Len(t, d.GetElements(), 0)
	})

	t.Run("MarkDownDocument", func(t *testing.T) {
		var d Document = markdowndocument.NewMarkDownDocument()

		require.Len(t, d.GetElements(), 0)
	})

}

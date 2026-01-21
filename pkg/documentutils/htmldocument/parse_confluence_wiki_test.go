package htmldocument_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/basicdocument"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/htmldocument"
)

func Test_ConfluenceWiki_SimplePageWithToc(t *testing.T) {

	rawHTML := `<html>
  <head></head>
  <body>
    <p>Text before table of contents.</p>
    <p>
      <ac:structured-macro ac:name="toc" ac:schema-version="1" ac:macro-id="85191b06-5f19-4ec1-94bd-541c7e1e0c4f"></ac:structured-macro>
    </p>
    <p>Text after table of contents.</p>
    <h1>Title</h1>
    <p>A simple text</p>
    <h2>Subtitle1</h2>
    <p>A simple text under subtitle 1</p>
    <h2>Subtitle2</h2>
    <p>A simple text under subtitle 2</p>
    <h1>Another Title</h1>
    <p>And this is the last example text of this page.</p>
  </body>
</html>`
	document, err := htmldocument.ParseString(rawHTML)
	require.NoError(t, err)
	require.NotNil(t, document)

	elements := document.GetElements()
	require.Len(t, elements, 10)

	text, err := elements[0].GetPlainText()
	require.NoError(t, err)
	require.EqualValues(t, "Text before table of contents.", text)
	_, ok := elements[0].(*basicdocument.Text)
	require.True(t, ok)

	text, err = elements[1].GetPlainText()
	require.NoError(t, err)
	require.EqualValues(t, "Text after table of contents.", text)
	_, ok = elements[1].(*basicdocument.Text)
	require.True(t, ok)

	text, err = elements[2].GetPlainText()
	require.NoError(t, err)
	require.EqualValues(t, "Title", text)
	_, ok = elements[2].(*basicdocument.Title)
	require.True(t, ok)
}

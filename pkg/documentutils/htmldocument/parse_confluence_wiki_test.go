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

func Test_ConfluenceWiki_CodeBlocks(t *testing.T) {
	rawHTML := `<html>
  <head></head>
  <body>
    <p>This is a code example without a language specified</p>
    <ac:structured-macro ac:name="code" ac:schema-version="1" ac:macro-id="12345678-abcd-1234-abcd-12345678910">
      <ac:plain-text-body>
        <!--[CDATA[hello world]]--></ac:plain-text-body>
    </ac:structured-macro>
    <p>And here Shell is explicitly specified:</p>
    <ac:structured-macro ac:name="code" ac:schema-version="1" ac:macro-id="12345678-abcd-1234-abcd-12345678911">
      <ac:parameter ac:name="language">shell</ac:parameter>
      <ac:plain-text-body>
        <!--[CDATA[echo "hello world!"]]--></ac:plain-text-body>
    </ac:structured-macro>
    <p>
      <br/>
    </p>
  </body>
</html>`
	document, err := htmldocument.ParseString(rawHTML)
	require.NoError(t, err)
	require.NotNil(t, document)

	elements := document.GetElements()
	require.Len(t, elements, 4)

	text, err := elements[0].GetPlainText()
	require.NoError(t, err)
	require.EqualValues(t, "This is a code example without a language specified", text)
	_, ok := elements[0].(*basicdocument.Text)
	require.True(t, ok)

	text, err = elements[1].GetPlainText()
	require.NoError(t, err)
	require.EqualValues(t, "hello world", text)
	_, ok = elements[1].(*basicdocument.CodeBlock)
	require.True(t, ok)

	text, err = elements[2].GetPlainText()
	require.NoError(t, err)
	require.EqualValues(t, "And here Shell is explicitly specified:", text)
	_, ok = elements[2].(*basicdocument.Text)
	require.True(t, ok)

	text, err = elements[3].GetPlainText()
	require.NoError(t, err)
	require.EqualValues(t, "echo \"hello world!\"", text)
	_, ok = elements[3].(*basicdocument.CodeBlock)
	require.True(t, ok)
}

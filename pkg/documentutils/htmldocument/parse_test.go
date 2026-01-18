package htmldocument_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/basicdocument"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/htmldocument"
)

func Test_ParseEmptyString(t *testing.T) {
	tests := []struct {
		toParse string
	}{
		{""},
		{" "},
		{"\t"},
		{"\n"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("parese empty '%s'", tt.toParse), func(t *testing.T) {
			parsed, err := htmldocument.ParseString(tt.toParse)
			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func TestParseEmptyDocument(t *testing.T) {
	rawHTMLWithBody := `
<!DOCTYPE html>
<html>
<body>
</body>
</html>`

	rawHTMLWithoutBody := `
<!DOCTYPE html>
<html>
</html>`

	tests := []struct {
		name    string
		rawHtml string
	}{
		{"with empty body", rawHTMLWithBody},
		{"without body", rawHTMLWithoutBody},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			document, err := htmldocument.ParseString(tt.rawHtml)
			require.NoError(t, err)
			require.NotNil(t, document)

			elements := document.GetElements()
			require.Len(t, elements, 0)
		})
	}
}

func Test_ParseExample1(t *testing.T) {
	rawHTML := `
<!DOCTYPE html>
<html>
<body>
    <h1>Main Title</h1>
    <h2>First Subtitle</h2>
	<p>This is the first paragraph under the first subtitle.</p>
	
    <h2>Second Subtitle</h2>
    <p>This is the second paragraph under the second subtitle.</p>
</body>
</html>`

	document, err := htmldocument.ParseString(rawHTML)
	require.NoError(t, err)
	require.NotNil(t, document)

	elements := document.GetElements()
	require.Len(t, elements, 5)

	text, err := elements[0].GetPlainText()
	require.NoError(t, err)
	require.EqualValues(t, "Main Title", text)
	_, ok := elements[0].(*basicdocument.Title)
	require.True(t, ok)

	text, err = elements[1].GetPlainText()
	require.NoError(t, err)
	require.EqualValues(t, "First Subtitle", text)
	_, ok = elements[1].(*basicdocument.SubTitle)
	require.True(t, ok)

	text, err = elements[2].GetPlainText()
	require.NoError(t, err)
	require.EqualValues(t, "This is the first paragraph under the first subtitle.", text)
	_, ok = elements[2].(*basicdocument.Text)
	require.True(t, ok)

	text, err = elements[3].GetPlainText()
	require.NoError(t, err)
	require.EqualValues(t, "Second Subtitle", text)
	_, ok = elements[3].(*basicdocument.SubTitle)
	require.True(t, ok)

	text, err = elements[4].GetPlainText()
	require.NoError(t, err)
	require.EqualValues(t, "This is the second paragraph under the second subtitle.", text)
	_, ok = elements[4].(*basicdocument.Text)
	require.True(t, ok)
}

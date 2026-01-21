package documentconverter_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/documentconverter"
)

func Test_Convert_HTML2md(t *testing.T) {
	rawHtml := `<html>
<body>
<h1>This is the title</h1>
<p>And this a text.</p>
</body>
</html>
`

	md, err := documentconverter.HtmlStringToMdString(rawHtml)
	require.NoError(t, err)

	expected := `# This is the title

And this a text.
`

	require.EqualValues(t, expected, md)
}

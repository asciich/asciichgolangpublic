package htmldocument_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/basicdocument"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/htmldocument"
)

func TestRenderNil(t *testing.T) {
	rendered, err := htmldocument.RenderAsString(nil)
	require.Error(t, err)
	require.Empty(t, rendered)
}

func TestRenderEmptyDocument(t *testing.T) {
	expected := `<html>
<body>
</body>
</html>
`

	document := basicdocument.NewBasicDocument()
	rendered, err := htmldocument.RenderAsString(document)
	require.NoError(t, err)
	require.EqualValues(t, expected, rendered)
}

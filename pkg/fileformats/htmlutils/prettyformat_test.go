package htmlutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/htmlutils"
)

func TestPrettyFormat(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		input := `<html><body><h1>Hello, World!</h1></body></html>`
		expectedOutput := `<html>
  <head></head>
  <body>
    <h1>Hello, World!</h1>
  </body>
</html>
`
		result, err := htmlutils.PrettyFormat(input)
		require.NoError(t, err, "Expected no error")
		require.Equal(t, expectedOutput, result, "The output should match the expected formatted HTML")
	})

	t.Run("empty input", func(t *testing.T) {
		input := ``
		expectedOutput := `<html>
  <head></head>
  <body></body>
</html>
`

		result, err := htmlutils.PrettyFormat(input)
		require.NoError(t, err, "Expected no error")
		require.Equal(t, expectedOutput, result, "The output for empty input should match the expected minimal HTML structure")
	})

	t.Run("href only", func(t *testing.T) {
		input := `<a href="target.html">linktext</a>`
		expectedOutput := `<html>
  <head></head>
  <body>
    <a href="target.html">linktext</a>
  </body>
</html>
`

		result, err := htmlutils.PrettyFormat(input)
		require.NoError(t, err, "Expected no error")
		require.Equal(t, expectedOutput, result, "The output for empty input should match the expected minimal HTML structure")
	})
}

func TestReformatHTMLAhrefIntext(t *testing.T) {
	input := `<html><body><h1>Hello, World!</h1><p>this <a href="target.html">is a link</a></p></body></html>`
	expectedOutput := `<html>
  <head></head>
  <body>
    <h1>Hello, World!</h1>
    <p>this 
      <a href="target.html">is a link</a>
    </p>
  </body>
</html>
`
	result, err := htmlutils.PrettyFormat(input)
	require.NoError(t, err, "Expected no error")
	require.Equal(t, expectedOutput, result, "The output should match the expected formatted HTML")
}

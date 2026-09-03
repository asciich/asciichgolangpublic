package gotemplateutils_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/templateutils/gotemplateutils"
	// adjust import alias/path if needed
)

// --- RenderTemplateFromStringAsString ---------------------------------------

func Test_RenderTemplateFromStringAsString(t *testing.T) {
	t.Run("simple variable substitution", func(t *testing.T) {
		rendered, err := gotemplateutils.RenderTemplateFromStringAsString(
			"Hello {{ .name }}!",
			map[string]interface{}{"name": "World"},
		)
		require.NoError(t, err)
		require.EqualValues(t, "Hello World!", rendered)
	})

	t.Run("multiple variables", func(t *testing.T) {
		rendered, err := gotemplateutils.RenderTemplateFromStringAsString(
			"{{ .greeting }}, {{ .name }}!",
			map[string]interface{}{"greeting": "Hi", "name": "Alice"},
		)
		require.NoError(t, err)
		require.EqualValues(t, "Hi, Alice!", rendered)
	})

	t.Run("ToUpper func", func(t *testing.T) {
		rendered, err := gotemplateutils.RenderTemplateFromStringAsString(
			"{{ ToUpper .name }}",
			map[string]interface{}{"name": "abc"},
		)
		require.NoError(t, err)
		require.EqualValues(t, "ABC", rendered)
	})

	t.Run("ToLower func", func(t *testing.T) {
		rendered, err := gotemplateutils.RenderTemplateFromStringAsString(
			"{{ ToLower .name }}",
			map[string]interface{}{"name": "ABC"},
		)
		require.NoError(t, err)
		require.EqualValues(t, "abc", rendered)
	})

	t.Run("template without variables", func(t *testing.T) {
		rendered, err := gotemplateutils.RenderTemplateFromStringAsString(
			"no placeholders here",
			map[string]interface{}{"unused": "value"},
		)
		require.NoError(t, err)
		require.EqualValues(t, "no placeholders here", rendered)
	})

	t.Run("empty inputString returns error", func(t *testing.T) {
		rendered, err := gotemplateutils.RenderTemplateFromStringAsString(
			"",
			map[string]interface{}{"name": "World"},
		)
		require.Error(t, err)
		require.EqualValues(t, "", rendered)
	})

	t.Run("nil variables returns error", func(t *testing.T) {
		rendered, err := gotemplateutils.RenderTemplateFromStringAsString(
			"Hello {{ .name }}!",
			nil,
		)
		require.Error(t, err)
		require.EqualValues(t, "", rendered)
	})

	t.Run("missing key returns error (missingkey=error)", func(t *testing.T) {
		rendered, err := gotemplateutils.RenderTemplateFromStringAsString(
			"Hello {{ .missing }}!",
			map[string]interface{}{"name": "World"},
		)
		require.Error(t, err)
		require.EqualValues(t, "", rendered)
	})

	t.Run("invalid template syntax returns error", func(t *testing.T) {
		rendered, err := gotemplateutils.RenderTemplateFromStringAsString(
			"Hello {{ .name ",
			map[string]interface{}{"name": "World"},
		)
		require.Error(t, err)
		require.EqualValues(t, "", rendered)
	})
}

// --- RenderTemplateFromFilePathAsString -------------------------------------

func Test_RenderTemplateFromFilePathAsString(t *testing.T) {
	ctx := context.Background()

	t.Run("renders template read from file path", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "template.txt")
		require.NoError(t, os.WriteFile(filePath, []byte("Hello {{ .name }}!"), 0o600))

		rendered, err := gotemplateutils.RenderTemplateFromFilePathAsString(
			ctx,
			filePath,
			map[string]interface{}{"name": "World"},
		)
		require.NoError(t, err)
		require.EqualValues(t, "Hello World!", rendered)
	})

	t.Run("empty inputFilePath returns error", func(t *testing.T) {
		rendered, err := gotemplateutils.RenderTemplateFromFilePathAsString(
			ctx,
			"",
			map[string]interface{}{"name": "World"},
		)
		require.Error(t, err)
		require.EqualValues(t, "", rendered)
	})

	t.Run("nil variables returns error", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "template.txt")
		require.NoError(t, os.WriteFile(filePath, []byte("Hello {{ .name }}!"), 0o600))

		rendered, err := gotemplateutils.RenderTemplateFromFilePathAsString(
			ctx,
			filePath,
			nil,
		)
		require.Error(t, err)
		require.EqualValues(t, "", rendered)
	})
}

// --- RenderTemplateFromFileAsString -----------------------------------------

func Test_RenderTemplateFromFileAsString(t *testing.T) {
	ctx := context.Background()

	t.Run("nil inputFile returns error", func(t *testing.T) {
		rendered, err := gotemplateutils.RenderTemplateFromFileAsString(
			ctx,
			nil,
			map[string]interface{}{"name": "World"},
		)
		require.Error(t, err)
		require.EqualValues(t, "", rendered)
	})
}

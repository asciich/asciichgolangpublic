package yamlutils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

// This test shows how to check if a file contains valid YAML.
func Test_Example_IsYamlFile(t *testing.T) {
	t.Run("simple yaml", func(t *testing.T) {
		// context used for testing:
		ctx := getCtx()

		// This simple YAML passes the validation:
		content := "---\n"
		content += "a: 123\n"

		// Write to temporary file:
		tempfile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, content)
		require.NoError(t, err)
		defer func() { _ = os.Remove(tempfile) }()

		// Check if file is a yaml file
		isYaml, err := yamlutils.IsYamlFile(ctx, tempfile, &yamlutils.ValidateOptions{})
		require.NoError(t, err)
		require.True(t, isYaml)
	})

	t.Run("empty string", func(t *testing.T) {
		// context used for testing:
		ctx := getCtx()

		// An empty string is not considered as valid YAML:
		content := ""

		// Write to temporary file:
		tempfile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, content)
		require.NoError(t, err)
		defer func() { _ = os.Remove(tempfile) }()

		// Check if file is a yaml file
		isYaml, err := yamlutils.IsYamlFile(ctx, tempfile, &yamlutils.ValidateOptions{})
		require.NoError(t, err)
		require.False(t, isYaml)
	})

	t.Run("invalid string", func(t *testing.T) {
		// context used for testing:
		ctx := getCtx()

		// This string is not valid YAML
		content := "a: b: Not yaml"

		// Write to temporary file:
		tempfile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, content)
		require.NoError(t, err)
		defer func() { _ = os.Remove(tempfile) }()

		// Check if file is a yaml file
		isYaml, err := yamlutils.IsYamlFile(ctx, tempfile, &yamlutils.ValidateOptions{})
		require.NoError(t, err)
		require.False(t, isYaml)
	})

	t.Run("simple json", func(t *testing.T) {
		// context used for testing:
		ctx := getCtx()

		// This simple JSON is a valid JSON.
		// Every valid JSON file is a YAML file as well by definition.
		content := "{\"a\": 123}"

		// Write to temporary file:
		tempfile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, content)
		require.NoError(t, err)
		defer func() { _ = os.Remove(tempfile) }()

		// Check if file is a yaml file
		isYaml, err := yamlutils.IsYamlFile(ctx, tempfile, &yamlutils.ValidateOptions{})
		require.NoError(t, err)
		require.True(t, isYaml) // a valid JSON file is a YAML file as well by definition.

		// If you want an error in case "only" a JSON document is provided as content:
		isYaml, err = yamlutils.IsYamlFile(ctx, tempfile, &yamlutils.ValidateOptions{
			RefuesePureJson: true,
		})
		require.NoError(t, err)
		require.False(t, isYaml) // Now the JSON file is not considered as valid YAML.
	})
}

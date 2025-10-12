package yamlutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/yamlutils"
)

// This test shows how to validate if a string contains yaml data or not.
func Test_Example_ValidateString(t *testing.T) {
	t.Run("simple yaml", func(t *testing.T) {
		// This simple YAML passes the validation:
		content := "---\n"
		content += "a: 123\n"

		// perform the validation:
		err := yamlutils.Validate(content, &yamlutils.ValidateOptions{})

		// No error indicates valid YAML:
		require.NoError(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		// An empty string is not considered as valid YAML:
		content := ""

		// perform the validation:
		err := yamlutils.Validate(content, &yamlutils.ValidateOptions{})

		// The error indicates no valid YAML:
		require.Error(t,err)
		require.ErrorIs(t, err, yamlutils.ErrInvalidYamlEmptyString)
	})

	t.Run("invalid string", func(t *testing.T) {
		// This string is not valid YAML
		content := "a: b: Not yaml"

		// perform the validation:
		err := yamlutils.Validate(content, &yamlutils.ValidateOptions{})

		// The error indicates no valid YAML:
		require.Error(t,err)
		require.ErrorIs(t, err, yamlutils.ErrInvalidYaml)
	})

	t.Run("simple json", func(t *testing.T) {
		// This simple JSON is a valid JSON.
		// Every valid JSON file is a YAML file as well by definition.
		content := "{\"a\": 123}"

		// perform the validation
		err := yamlutils.Validate(content, &yamlutils.ValidateOptions{})

		// No error indicates valid YAML:
		require.NoError(t, err)

		// If you want an error in case "only" a JSON document is provided as content:
		err = yamlutils.Validate(content, &yamlutils.ValidateOptions{
			RefuesePureJson: true, // Explicitly refuse pure JSON documents as valid YAML
		})

		// The error indicates not a valid YAML:
		require.Error(t, err)
		require.ErrorIs(t, err, yamlutils.ErrOnlyJSONinDocument)
	})

}

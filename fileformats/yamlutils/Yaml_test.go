package yamlutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestYaml_runYqQueryAgainstYamlStringAsString_2(t *testing.T) {
	resourceName := "resourceName"
	namespaceName := "namespaceName"

	roleYaml := ""
	roleYaml += "apiVersion: v1\n"
	roleYaml += "kind: Secret\n"
	roleYaml += "metadata:\n"
	roleYaml += "  name: " + resourceName + "\n"
	roleYaml += "  namespace: " + namespaceName + "\n"

	expectedYaml := ""
	expectedYaml += "apiVersion: v1\n"
	expectedYaml += "kind: hello_world\n"
	expectedYaml += "metadata:\n"
	expectedYaml += "  name: " + resourceName + "\n"
	expectedYaml += "  namespace: " + namespaceName

	assert.EqualValues(
		t,
		expectedYaml,
		MustRunYqQueryAginstYamlStringAsString(roleYaml, ".kind=\"hello_world\""),
	)
}

func TestYaml_runYqQueryAgainstYamlStringAsString(t *testing.T) {
	tests := []struct {
		yamlString     string
		query          string
		expectedResult string
	}{
		{"---\na: 1234\nb: 456", ".a", "1234"},
		{"---\na: 1234\nb: 456", ".a=123", "---\na: 123\nb: 456"},
		{"a: 1234\nb: 456", ".a=123", "a: 123\nb: 456"},
		{"a:\n  b: 1\n  c: 2", ".a.b", "1"},
		{"a:\n  b: 1\n  c: 2", ".a.c", "2"},
		{"a:\n  b: 1\n  c: 2", ".a.b=3", "a:\n  b: 3\n  c: 2"},
		{"a:\n  b: 1\n  c: 2", ".a.c=3", "a:\n  b: 1\n  c: 3"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedResult,
					MustRunYqQueryAginstYamlStringAsString(tt.yamlString, tt.query),
				)
			},
		)
	}
}

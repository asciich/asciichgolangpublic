package yamlutils_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

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

	got, err := yamlutils.RunYqQueryAginstYamlStringAsString(roleYaml, ".kind=\"hello_world\"")
	require.NoError(t, err)
	require.EqualValues(t, expectedYaml, got)
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
				got, err := yamlutils.RunYqQueryAginstYamlStringAsString(tt.yamlString, tt.query)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedResult, got)
			},
		)
	}
}

func Test_SplitMultiYaml(t *testing.T) {
	t.Run("split single", func(t *testing.T) {
		content := "a: 5"
		require.EqualValues(
			t,
			[]string{"---\n" + content + "\n"},
			yamlutils.SplitMultiYaml(content),
		)
	})

	t.Run("split single indent", func(t *testing.T) {
		content := "a:\n  b: 5"
		require.EqualValues(
			t,
			[]string{"---\n" + content + "\n"},
			yamlutils.SplitMultiYaml(content),
		)
	})

	t.Run("split single and doxument start", func(t *testing.T) {
		content := "a: 5"
		require.EqualValues(
			t,
			[]string{"---\n" + content + "\n"},
			yamlutils.SplitMultiYaml("---\n"+content),
		)
	})

	t.Run("split double and doxument start", func(t *testing.T) {
		content := "a: 5"
		content2 := "b: 1"
		require.EqualValues(
			t,
			[]string{"---\n" + content + "\n", "---\n" + content2 + "\n"},
			yamlutils.SplitMultiYaml("---\n"+content+"\n---\n"+content2),
		)
	})

	t.Run("split double and multiple document start", func(t *testing.T) {
		content := "a: 5"
		content2 := "b: 1"
		require.EqualValues(
			t,
			[]string{"---\n" + content + "\n", "---\n" + content2 + "\n"},
			yamlutils.SplitMultiYaml("---\n"+content+"\n---\n---\n---\n"+content2),
		)
	})

	t.Run("split double and commented out", func(t *testing.T) {
		content := "a: 5"
		content2 := "b: 1"
		require.EqualValues(
			t,
			[]string{"---\na: 5\n#---\n#c: 16 is commented out\n", "---\n" + content2 + "\n"},
			yamlutils.SplitMultiYaml("---\n"+content+"\n#---\n#c: 16 is commented out\n---\n---\n"+content2),
		)
	})
}

func Test_MergeMultiYaml(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		require.EqualValues(t, "\n", yamlutils.MustMergeMultiYaml([]string{""}))
		require.EqualValues(t, "\n", yamlutils.MustMergeMultiYaml([]string{" "}))
		require.EqualValues(t, "\n", yamlutils.MustMergeMultiYaml([]string{"\n"}))
	})

	t.Run("single entry", func(t *testing.T) {
		require.EqualValues(t, "---\na: 1234\n", yamlutils.MustMergeMultiYaml([]string{"a: 1234"}))
	})

	t.Run("single entry leading yaml start marker", func(t *testing.T) {
		require.EqualValues(t, "---\na: 1234\n", yamlutils.MustMergeMultiYaml([]string{"---\na: 1234"}))
	})

	t.Run("double entry", func(t *testing.T) {
		require.EqualValues(t, "---\na: 1234\n---\nb: 123\n", yamlutils.MustMergeMultiYaml([]string{"a: 1234", "b: 123"}))
		require.EqualValues(t, "---\na: 1234\n---\nb: 123\n", yamlutils.MustMergeMultiYaml([]string{"a: 1234", "---\nb: 123"}))
		require.EqualValues(t, "---\na: 1234\n---\nb: 123\n", yamlutils.MustMergeMultiYaml([]string{"---\na: 1234", "b: 123"}))
		require.EqualValues(t, "---\na: 1234\n---\nb: 123\n", yamlutils.MustMergeMultiYaml([]string{"---\na: 1234", "---\nb: 123"}))
		require.EqualValues(t, "---\na: 1234\n---\nb: 123\n", yamlutils.MustMergeMultiYaml([]string{"---\na: 1234\n", "---\nb: 123\n"}))
	})

	t.Run("double entry leading whitespces", func(t *testing.T) {
		require.EqualValues(t, "---\na: 1234\n---\nb: 123\n", yamlutils.MustMergeMultiYaml([]string{"\n\n\n---\na: 1234\n", "\n\n\n---\nb: 123\n"}))
	})

	t.Run("single intend", func(t *testing.T) {
		require.EqualValues(
			t,
			"---\na:\n  b: 1\n",
			yamlutils.MustMergeMultiYaml([]string{"a:\n  b: 1\n"}),
		)
	})
}

func Test_Validate(t *testing.T) {
	for i, s := range []string{"", " ", "\n", "\t", "\n\t     "} {
		t.Run(
			"empty string"+strconv.Itoa(i),
			func(t *testing.T) {
				err := yamlutils.Validate(s, &yamlutils.ValidateOptions{})
				require.ErrorIs(t, err, yamlutils.ErrInvalidYaml)
				require.ErrorIs(t, err, yamlutils.ErrInvalidYamlEmptyString)
			},
		)
	}

	t.Run("single char", func(t *testing.T) {
		require.NoError(t, yamlutils.Validate("x", &yamlutils.ValidateOptions{}))
	})

	t.Run("Only document start", func(t *testing.T) {
		require.NoError(t, yamlutils.Validate("---", &yamlutils.ValidateOptions{}))
		require.NoError(t, yamlutils.Validate("---\n", &yamlutils.ValidateOptions{}))
	})

	t.Run("single key value", func(t *testing.T) {
		require.NoError(t, yamlutils.Validate("a: b", &yamlutils.ValidateOptions{}))
		require.NoError(t, yamlutils.Validate("---\na: b", &yamlutils.ValidateOptions{}))
	})

	t.Run("key without value", func(t *testing.T) {
		err := yamlutils.Validate("a: b: 5\n", &yamlutils.ValidateOptions{})

		require.ErrorIs(t, err, yamlutils.ErrInvalidYaml)
		require.NotErrorIs(t, err, yamlutils.ErrInvalidYamlEmptyString)
	})
}

func Test_EnsureDocumentStart(t *testing.T) {
	t.Run("Empty string", func(t *testing.T) {
		require.EqualValues(t, "---\n", yamlutils.EnsureDocumentStart(""))
	})

	t.Run("document only", func(t *testing.T) {
		require.EqualValues(t, "---\n", yamlutils.EnsureDocumentStart("---"))
	})

	t.Run("document and newline", func(t *testing.T) {
		require.EqualValues(t, "---\n", yamlutils.EnsureDocumentStart("---\n"))
	})

	t.Run("key value", func(t *testing.T) {
		require.EqualValues(t, "---\na: b", yamlutils.EnsureDocumentStart("a: b"))
	})

	t.Run("comment only", func(t *testing.T) {
		require.EqualValues(t, "---\n# comment", yamlutils.EnsureDocumentStart("# comment"))
	})

	t.Run("comment only and key value", func(t *testing.T) {
		require.EqualValues(t, "---\n# comment\na: b", yamlutils.EnsureDocumentStart("# comment\na: b"))
	})

	t.Run("comment only and document start only", func(t *testing.T) {
		require.EqualValues(t, "# comment\n---\n", yamlutils.EnsureDocumentStart("# comment\n---"))
	})
}

func Test_EnsureDocumentStartAndEnd(t *testing.T) {
	t.Run("key value", func(t *testing.T) {
		require.EqualValues(t, "---\na: 42\n", yamlutils.EnsureDocumentStartAndEnd("a: 42"))
	})
}

func Test_IsYamlString(t *testing.T) {
	t.Run("empty", func(t *testing.T) {

	})
}

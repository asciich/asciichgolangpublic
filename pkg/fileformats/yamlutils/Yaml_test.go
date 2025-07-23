package yamlutils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/testutils"
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

	require.EqualValues(
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
				require := require.New(t)

				require.EqualValues(
					tt.expectedResult,
					MustRunYqQueryAginstYamlStringAsString(tt.yamlString, tt.query),
				)
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
			SplitMultiYaml(content),
		)
	})

	t.Run("split single indent", func(t *testing.T) {
		content := "a:\n  b: 5"
		require.EqualValues(
			t,
			[]string{"---\n" + content + "\n"},
			SplitMultiYaml(content),
		)
	})

	t.Run("split single and doxument start", func(t *testing.T) {
		content := "a: 5"
		require.EqualValues(
			t,
			[]string{"---\n" + content + "\n"},
			SplitMultiYaml("---\n"+content),
		)
	})

	t.Run("split double and doxument start", func(t *testing.T) {
		content := "a: 5"
		content2 := "b: 1"
		require.EqualValues(
			t,
			[]string{"---\n" + content + "\n", "---\n" + content2 + "\n"},
			SplitMultiYaml("---\n"+content+"\n---\n"+content2),
		)
	})

	t.Run("split double and multiple document start", func(t *testing.T) {
		content := "a: 5"
		content2 := "b: 1"
		require.EqualValues(
			t,
			[]string{"---\n" + content + "\n", "---\n" + content2 + "\n"},
			SplitMultiYaml("---\n"+content+"\n---\n---\n---\n"+content2),
		)
	})

	t.Run("split double and commented out", func(t *testing.T) {
		content := "a: 5"
		content2 := "b: 1"
		require.EqualValues(
			t,
			[]string{"---\na: 5\n#---\n#c: 16 is commented out\n", "---\n" + content2 + "\n"},
			SplitMultiYaml("---\n"+content+"\n#---\n#c: 16 is commented out\n---\n---\n"+content2),
		)
	})
}

func Test_MergeMultiYaml(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		require.EqualValues(t, "\n", MustMergeMultiYaml([]string{""}))
		require.EqualValues(t, "\n", MustMergeMultiYaml([]string{" "}))
		require.EqualValues(t, "\n", MustMergeMultiYaml([]string{"\n"}))
	})

	t.Run("single entry", func(t *testing.T) {
		require.EqualValues(t, "---\na: 1234\n", MustMergeMultiYaml([]string{"a: 1234"}))
	})

	t.Run("single entry leading yaml start marker", func(t *testing.T) {
		require.EqualValues(t, "---\na: 1234\n", MustMergeMultiYaml([]string{"---\na: 1234"}))
	})

	t.Run("double entry", func(t *testing.T) {
		require.EqualValues(t, "---\na: 1234\n---\nb: 123\n", MustMergeMultiYaml([]string{"a: 1234", "b: 123"}))
		require.EqualValues(t, "---\na: 1234\n---\nb: 123\n", MustMergeMultiYaml([]string{"a: 1234", "---\nb: 123"}))
		require.EqualValues(t, "---\na: 1234\n---\nb: 123\n", MustMergeMultiYaml([]string{"---\na: 1234", "b: 123"}))
		require.EqualValues(t, "---\na: 1234\n---\nb: 123\n", MustMergeMultiYaml([]string{"---\na: 1234", "---\nb: 123"}))
		require.EqualValues(t, "---\na: 1234\n---\nb: 123\n", MustMergeMultiYaml([]string{"---\na: 1234\n", "---\nb: 123\n"}))
	})

	t.Run("double entry leading whitespces", func(t *testing.T) {
		require.EqualValues(t, "---\na: 1234\n---\nb: 123\n", MustMergeMultiYaml([]string{"\n\n\n---\na: 1234\n", "\n\n\n---\nb: 123\n"}))
	})

	t.Run("single intend", func(t *testing.T) {
		require.EqualValues(
			t,
			"---\na:\n  b: 1\n",
			MustMergeMultiYaml([]string{"a:\n  b: 1\n"}),
		)
	})
}

func Test_Validate(t *testing.T) {
	for i, s := range []string{"", " ", "\n", "\t", "\n\t     "} {
		i := i
		s := s
		t.Run(
			"empty string"+strconv.Itoa(i),
			func(t *testing.T) {
				err := Validate(s)
				require.ErrorIs(t, err, ErrInvalidYaml)
				require.ErrorIs(t, err, ErrInvalidYamlEmptyString)
			},
		)
	}

	t.Run("single char", func(t *testing.T) {
		require.NoError(t, Validate("x"))
	})

	t.Run("Only document start", func(t *testing.T) {
		require.NoError(t, Validate("---"))
		require.NoError(t, Validate("---\n"))
	})

	t.Run("single key value", func(t *testing.T) {
		require.NoError(t, Validate("a: b"))
		require.NoError(t, Validate("---\na: b"))
	})

	t.Run("key without value", func(t *testing.T) {
		err := Validate("a: b: 5\n")

		require.ErrorIs(t, err, ErrInvalidYaml)
		require.NotErrorIs(t, err, ErrInvalidYamlEmptyString)
	})
}

func Test_EnsureDocumentStart(t *testing.T) {
	t.Run("Empty string", func(t *testing.T) {
		require.EqualValues(t, "---\n", EnsureDocumentStart(""))
	})

	t.Run("document only", func(t *testing.T) {
		require.EqualValues(t, "---\n", EnsureDocumentStart("---"))
	})

	t.Run("document and newline", func(t *testing.T) {
		require.EqualValues(t, "---\n", EnsureDocumentStart("---\n"))
	})

	t.Run("key value", func(t *testing.T) {
		require.EqualValues(t, "---\na: b", EnsureDocumentStart("a: b"))
	})

	t.Run("comment only", func(t *testing.T) {
		require.EqualValues(t, "---\n# comment", EnsureDocumentStart("# comment"))
	})

	t.Run("comment only and key value", func(t *testing.T) {
		require.EqualValues(t, "---\n# comment\na: b", EnsureDocumentStart("# comment\na: b"))
	})

	t.Run("comment only and document start only", func(t *testing.T) {
		require.EqualValues(t, "# comment\n---\n", EnsureDocumentStart("# comment\n---"))
	})
}

func Test_EnsureDocumentStartAndEnd(t *testing.T) {
	t.Run("key value", func(t *testing.T) {
		require.EqualValues(t, "---\na: 42\n", EnsureDocumentStartAndEnd("a: 42"))
	})
}

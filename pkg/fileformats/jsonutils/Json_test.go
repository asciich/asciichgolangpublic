package jsonutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestJsonRunJqAgainstJsonStringAsString(t *testing.T) {
	tests := []struct {
		jsonString     string
		query          string
		expectedResult string
	}{
		{"{\"a\": 15}", ".a", "15"},
		{"{\"a\": 15, \"b\": 16}", ".", "{\n    \"a\": 15,\n    \"b\": 16\n}"},
		{"{\"a\": 15, \"b\": 16}", ".a", "15"},
		{"{\"a\": 15, \"b\": 16}", ".b", "16"},
		{"{\"a\": 15, \"b\": 16}", "del(.b)", "{\n    \"a\": 15\n}"},
		{"{\"a\": 15, \"hello\": \"world\"}", ".hello", "world"},
		{"{\"a\": 15, \"b\": {\"c\": 13, \"d\": \"efg\"} }", ".b", "{\n    \"c\": 13,\n    \"d\": \"efg\"\n}"},
		{"{\"a\": 15, \"b\": [\"c\", \"d\", \"efg\"] }", ".b", "[\n    \"c\",\n    \"d\",\n    \"efg\"\n]"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				result := MustRunJqAgainstJsonStringAsString(tt.jsonString, tt.query)

				require.EqualValues(tt.expectedResult, result)
			},
		)
	}
}

func TestJsonLoadKeyValueDict(t *testing.T) {
	tests := []struct {
		jsonString     string
		expectedResult map[string]string
	}{
		{"{}", map[string]string{}},
		{"{\"a\": 15}", map[string]string{"a": "15"}},
		{"{\"a\": 15, \"hello\": \"world\"}", map[string]string{"a": "15", "hello": "world"}},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				result := MustLoadKeyValueStringDictFromJsonString(tt.jsonString)

				require.EqualValues(tt.expectedResult, result)
			},
		)
	}
}

func TestJsonPrettyFormatJsonString(t *testing.T) {
	tests := []struct {
		jsonString     string
		expectedResult string
	}{
		{"{}", "{}\n"},
		{"{\"a\": 15}", "{\n    \"a\": 15\n}\n"},
		{"{\"a\": 15, \"hello\": \"world\"}", "{\n    \"a\": 15,\n    \"hello\": \"world\"\n}\n"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				result := MustPrettyFormatJsonString(tt.jsonString)

				require.EqualValues(tt.expectedResult, result)
			},
		)
	}
}

func TestJsonStringToYamlString(t *testing.T) {
	tests := []struct {
		jsonString     string
		expectedResult string
	}{
		{"{\"a\": 15}", "---\na: 15\n"},
		{"{\"a\": 15, \"hello\": \"world\"}", "---\na: 15\nhello: world\n"},
		{"{\"a\": 15, \"b\" : {\"hello\": \"world\"} }", "---\na: 15\nb:\n    hello: world\n"},
		{"[1,2,3]", "---\n- 1\n- 2\n- 3\n"},
		{"{\"a\": [1,2,3]}", "---\na:\n    - 1\n    - 2\n    - 3\n"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				result := MustJsonStringToYamlString(tt.jsonString)

				require.EqualValues(tt.expectedResult, result)
			},
		)
	}
}

func TestJsonStringToYamlFileByPath(t *testing.T) {
	tests := []struct {
		jsonString     string
		expectedResult string
	}{
		{"{\"a\": 15}", "---\na: 15\n"},
		{"{\"a\": 15, \"hello\": \"world\"}", "---\na: 15\nhello: world\n"},
		{"{\"a\": 15, \"b\" : {\"hello\": \"world\"} }", "---\na: 15\nb:\n    hello: world\n"},
		{"[1,2,3]", "---\n- 1\n- 2\n- 3\n"},
		{"{\"a\": [1,2,3]}", "---\na:\n    - 1\n    - 2\n    - 3\n"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				emptyFile := tempfiles.MustCreateEmptyTemporaryFile(verbose)
				defer emptyFile.MustDelete(verbose)

				createdFile := MustJsonStringToYamlFileByPath(tt.jsonString, emptyFile.MustGetLocalPath(), verbose)

				require.EqualValues(tt.expectedResult, createdFile.MustReadAsString())
			},
		)
	}
}

func TestJsonStringHas(t *testing.T) {
	tests := []struct {
		jsonString     string
		query          string
		keyToCheck     string
		expectedResult bool
	}{
		{"{\"a\": 15}", ".", "a", true},
		{"{\"a\": 15}", ".", "b", false},
		{"{\"b\": 15}", ".", "a", false},
		{"{\"b\": 15}", ".", "b", true},
		{"{\"b\": { \"c\": 123 }}", ".", "a", false},
		{"{\"b\": { \"c\": 123 }}", ".", "b", true},
		{"{\"b\": { \"c\": 123 }}", ".b", "a", false},
		{"{\"b\": { \"c\": 123 }}", ".b", "b", false},
		{"{\"b\": { \"c\": 123 }}", ".b", "c", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				require.EqualValues(
					tt.expectedResult,
					MustJsonStringHas(tt.jsonString, tt.query, tt.keyToCheck),
				)
			},
		)
	}
}

func TestJsonFileHas(t *testing.T) {
	tests := []struct {
		jsonString     string
		query          string
		keyToCheck     string
		expectedResult bool
	}{
		{"{\"a\": 15}", ".", "a", true},
		{"{\"a\": 15}", ".", "b", false},
		{"{\"b\": 15}", ".", "a", false},
		{"{\"b\": 15}", ".", "b", true},
		{"{\"b\": { \"c\": 123 }}", ".", "a", false},
		{"{\"b\": { \"c\": 123 }}", ".", "b", true},
		{"{\"b\": { \"c\": 123 }}", ".b", "a", false},
		{"{\"b\": { \"c\": 123 }}", ".b", "b", false},
		{"{\"b\": { \"c\": 123 }}", ".b", "c", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				tempFile := tempfiles.MustCreateFromString(tt.jsonString, verbose)
				defer tempFile.Delete(verbose)

				require.EqualValues(
					tt.expectedResult,
					MustJsonFileByPathHas(tempFile.MustGetLocalPath(), tt.query, tt.keyToCheck),
				)
			},
		)
	}
}

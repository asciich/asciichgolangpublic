package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				result := Json().MustRunJqAgainstJsonStringAsString(tt.jsonString, tt.query)

				assert.EqualValues(tt.expectedResult, result)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				result := Json().MustLoadKeyValueStringDictFromJsonString(tt.jsonString)

				assert.EqualValues(tt.expectedResult, result)
			},
		)
	}
}

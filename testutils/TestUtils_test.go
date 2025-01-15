package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestsFormatAsTestname(t *testing.T) {
	tests := []struct {
		objectToFormat   interface{}
		expectedTestname string
	}{
		{"only a string", "only_a_string"},
		{struct{ a string }{a: "hello"}, "hello"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				testname := MustFormatAsTestname(tt.objectToFormat)

				assert.EqualValues(tt.expectedTestname, testname)
			},
		)
	}
}

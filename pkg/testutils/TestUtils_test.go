package testutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				testname, err := testutils.FormatAsTestname(tt.objectToFormat)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedTestname, testname)
			},
		)
	}
}

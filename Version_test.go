package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestVersionEquals(t *testing.T) {
	tests := []struct {
		version1       string
		version2       string
		expectedEquals bool
	}{
		{"20231231_203102", "20231231_203102", true},
		{"20231231_203102", "20231231_203103", false},
		{"20231231_203102", "v0.1.2", false},
		{"20231231_203102", "0.1.2", false},
		{"0.1.2", "0.1.2", true},
		{"v0.1.2", "0.1.2", true},
		{"v0.1.2", "V0.1.2", true},
		{"v0.1.2", "V0.1.3", false},
		{"v0.1.2", "V0.3.2", false},
		{"v1.1.2", "V0.1.2", false},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				version1 := Versions().MustGetNewVersionByString(tt.version1)
				version2 := Versions().MustGetNewVersionByString(tt.version2)

				require.EqualValues(
					tt.expectedEquals,
					version1.Equals(version2),
				)
			},
		)
	}
}

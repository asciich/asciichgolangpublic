package versionutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestVersionSemanticVersionGetNextVersion(t *testing.T) {
	tests := []struct {
		versionString       string
		nextVersionType     string
		expectedNextVersion string
	}{
		{"0.0.1", "patch", "v0.0.2"},
		{"0.0.1", "minor", "v0.1.0"},
		{"0.0.1", "major", "v1.0.0"},
		{"v1.2.3", "patch", "v1.2.4"},
		{"v1.2.3", "minor", "v1.3.0"},
		{"v1.2.3", "major", "v2.0.0"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				var version Version = Versions().MustGetNewVersionByString(tt.versionString)
				nextVersion := version.MustGetNextVersion(tt.nextVersionType)

				require.EqualValues(
					tt.expectedNextVersion,
					nextVersion.MustGetAsString(),
				)
			},
		)
	}
}

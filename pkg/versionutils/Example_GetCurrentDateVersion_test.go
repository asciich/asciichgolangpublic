package versionutils_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/versionutils"
)

// Demonstrates how to get a new DateVersion with the current date and time set.
func Test_Example_GetCurrentDateVersion(t *testing.T) {
	// Get current date and time as new date version:
	dateVersion := versionutils.NewCurrentDateVersion()

	// Get date version as string
	versionString, err := dateVersion.GetAsString()
	require.NoError(t, err)

	// Check date version string matches expected format of "YYYYmmdd_HHMMSS":
	require.Regexp(t, regexp.MustCompile(`^20[0-9]{6}_[0-9]{6}$`), versionString)
}

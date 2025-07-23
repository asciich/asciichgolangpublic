package versionutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/versionutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/testutils"
)

func TestVersions_GetDateVersionString(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		dateVersion := versionutils.NewCurrentDateVersion()

		dateVersionString, err := dateVersion.GetAsString()
		require.NoError(t, err)

		require.Len(t, dateVersionString, len("YYYYmmdd_HHMMSS"))
		require.NoError(t, versionutils.CheckIsDateVersionString(dateVersionString))
	})
}

func TestVersions_GetSoftwareVersionEnvVarName(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		require := require.New(t)

		require.EqualValues("SOFTWARE_VERSION", versionutils.GetSoftwareVersionEnvVarName())
	})
}

func TestVersions_IsDateVersionString(t *testing.T) {
	tests := []struct {
		versionString           string
		expectedIsVersionString bool
	}{
		{"testcase", false},
		{"20231112_123456", true},
		{"20231112_12345", false},
		{"20231112_1234566", false},
		{"20231112_a23456", false},
		{"v20231112_123456", true},
		{"v20231112_12345", false},
		{"v20231112_1234566", false},
		{"v20231112_a23456", false},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expectedIsVersionString, versionutils.IsDateVersionString(tt.versionString))
			},
		)
	}
}

func TestVersions_IsVersionString(t *testing.T) {
	tests := []struct {
		versionString           string
		expectedIsVersionString bool
	}{
		{"testcase", false},
		{"20231112_123456", true},
		{"20231112_12345", false},
		{"20231112_1234566", false},
		{"20231112_a23456", false},
		{"0.0.1", true},
		{"0.0.10", true},
		{"0.0.100", true},
		{"0.1.0", true},
		{"0.10.0", true},
		{"0.100.0", true},
		{"2.100.0", true},
		{"20.100.0", true},
		{"200.100.0", true},
		{"200.100.3", true},
		{"200.100.32", true},
		{"200.100.320", true},
		{"v0.0.1", true},
		{"v0.0.10", true},
		{"v0.0.100", true},
		{"v0.1.0", true},
		{"v0.10.0", true},
		{"v0.100.0", true},
		{"v2.100.0", true},
		{"v20.100.0", true},
		{"v200.100.0", true},
		{"v200.100.3", true},
		{"v200.100.32", true},
		{"v200.100.320", true},
		{"V0.0.1", true},
		{"V0.0.10", true},
		{"V0.0.100", true},
		{"V0.1.0", true},
		{"V0.10.0", true},
		{"V0.100.0", true},
		{"V2.100.0", true},
		{"V20.100.0", true},
		{"V200.100.0", true},
		{"V200.100.3", true},
		{"V200.100.32", true},
		{"V200.100.320", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expectedIsVersionString, versionutils.IsVersionString(tt.versionString))
			},
		)
	}
}

func TestVersions_IsSemanticVersion(t *testing.T) {
	tests := []struct {
		versionString             string
		expectedIsSemanticVersion bool
	}{
		{"20231112_123456", false},
		{"0.0.1", true},
		{"0.0.10", true},
		{"0.0.100", true},
		{"0.1.0", true},
		{"0.10.0", true},
		{"0.100.0", true},
		{"2.100.0", true},
		{"20.100.0", true},
		{"200.100.0", true},
		{"200.100.3", true},
		{"200.100.32", true},
		{"200.100.320", true},
		{"v0.0.1", true},
		{"v0.0.10", true},
		{"v0.0.100", true},
		{"v0.1.0", true},
		{"v0.10.0", true},
		{"v0.100.0", true},
		{"v2.100.0", true},
		{"v20.100.0", true},
		{"v200.100.0", true},
		{"v200.100.3", true},
		{"v200.100.32", true},
		{"v200.100.320", true},
		{"V0.0.1", true},
		{"V0.0.10", true},
		{"V0.0.100", true},
		{"V0.1.0", true},
		{"V0.10.0", true},
		{"V0.100.0", true},
		{"V2.100.0", true},
		{"V20.100.0", true},
		{"V200.100.0", true},
		{"V200.100.3", true},
		{"V200.100.32", true},
		{"V200.100.320", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				var version versionutils.Version
				var err error

				version, err = versionutils.ReadFromString(tt.versionString)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedIsSemanticVersion, version.IsSemanticVersion())
			},
		)
	}
}

func TestVersions_GetLatestVersionFromSlice(t *testing.T) {
	tests := []struct {
		versionStrings []string
		expectedNewest string
	}{
		{[]string{"v0.0.0"}, "v0.0.0"},
		{[]string{"v0.0.0", "v0.0.1"}, "v0.0.1"},

		{[]string{"v0.0.0", "v0.0.9", "v0.0.1"}, "v0.0.9"},
		{[]string{"v0.0.0", "v0.0.9", "v0.0.11"}, "v0.0.11"},
		{[]string{"v0.0.0", "v0.0.11", "v0.0.9"}, "v0.0.11"},

		{[]string{"v0.0.0", "v0.9.0", "v0.1.0"}, "v0.9.0"},
		{[]string{"v0.0.0", "v0.9.0", "v0.11.0"}, "v0.11.0"},
		{[]string{"v0.0.0", "v0.11.0", "v0.9.0"}, "v0.11.0"},

		{[]string{"v0.0.0", "v9.0.0", "v1.0.0"}, "v9.0.0"},
		{[]string{"v0.0.0", "v9.0.0", "v11.0.0"}, "v11.0.0"},
		{[]string{"v0.0.0", "v11.0.0", "v9.0.0"}, "v11.0.0"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				versions, err := versionutils.GetVersionsFromStringSlice(tt.versionStrings)
				require.NoError(t, err)

				latestVersion, err := versionutils.GetLatestVersionFromSlice(versions)
				require.NoError(t, err)

				expectedNewestVersion, err := versionutils.ReadFromString(tt.expectedNewest)
				require.NoError(t, err)

				require.True(t, latestVersion.Equals(expectedNewestVersion))
			},
		)
	}
}

func TestVersions_SortStringSlice(t *testing.T) {
	tests := []struct {
		versionStrings []string
		expectedSorted []string
	}{
		{[]string{"v0.0.0"}, []string{"v0.0.0"}},
		{[]string{"v0.0.0", "v0.1.2"}, []string{"v0.0.0", "v0.1.2"}},
		{[]string{"v0.1.2", "v0.0.0"}, []string{"v0.0.0", "v0.1.2"}},
	}
	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expectedSorted, mustutils.Must(versionutils.SortStringSlice(tt.versionStrings)))
			},
		)
	}
}

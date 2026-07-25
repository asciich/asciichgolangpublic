package versionutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

func TestVersions_GetDateVersionString(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		dateVersion := versionutils.NewCurrentDateVersion()

		dateVersionString, err := dateVersion.GetAsString()
		require.NoError(t, err)

		require.Len(t, dateVersionString, len("YYYYmmdd_HHMMSS"))
		require.NoError(t, versionutils.CheckDateVersionString(dateVersionString))
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

				version, err = versionutils.NewFromString(tt.versionString)
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

				expectedNewestVersion, err := versionutils.NewFromString(tt.expectedNewest)
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

func TestCheckDateVersionString(t *testing.T) {
	t.Run("valid date version string", func(t *testing.T) {
		require.NoError(t, versionutils.CheckDateVersionString("20260725_143059"))
	})

	t.Run("valid midnight", func(t *testing.T) {
		require.NoError(t, versionutils.CheckDateVersionString("20260101_000000"))
	})

	t.Run("valid end of day", func(t *testing.T) {
		require.NoError(t, versionutils.CheckDateVersionString("20261231_235959"))
	})

	t.Run("empty string", func(t *testing.T) {
		require.Error(t, versionutils.CheckDateVersionString(""))
	})

	t.Run("too short", func(t *testing.T) {
		require.Error(t, versionutils.CheckDateVersionString("20260725"))
	})

	t.Run("too long", func(t *testing.T) {
		require.Error(t, versionutils.CheckDateVersionString("20260725_1430591"))
	})

	t.Run("missing underscore separator", func(t *testing.T) {
		require.Error(t, versionutils.CheckDateVersionString("20260725-143059"))
	})

	t.Run("no separator at all", func(t *testing.T) {
		require.Error(t, versionutils.CheckDateVersionString("202607251430590"))
	})

	t.Run("letters in date part", func(t *testing.T) {
		require.Error(t, versionutils.CheckDateVersionString("2026ab25_143059"))
	})

	t.Run("letters in time part", func(t *testing.T) {
		require.Error(t, versionutils.CheckDateVersionString("20260725_14ab59"))
	})

	t.Run("spaces", func(t *testing.T) {
		require.Error(t, versionutils.CheckDateVersionString("20260725 143059"))
	})

	t.Run("semantic version is not a date version", func(t *testing.T) {
		require.Error(t, versionutils.CheckDateVersionString("1.2.3"))
	})

	t.Run("special characters", func(t *testing.T) {
		require.Error(t, versionutils.CheckDateVersionString("2026/07/25_14:30"))
	})
}

func TestCheckSemanticVersionString(t *testing.T) {
	t.Run("valid three part version", func(t *testing.T) {
		require.NoError(t, versionutils.CheckSemanticVersionString("1.2.3"))
	})

	t.Run("valid with zeros", func(t *testing.T) {
		require.NoError(t, versionutils.CheckSemanticVersionString("0.0.0"))
	})

	t.Run("valid high numbers", func(t *testing.T) {
		require.NoError(t, versionutils.CheckSemanticVersionString("100.200.300"))
	})

	t.Run("valid single digits", func(t *testing.T) {
		require.NoError(t, versionutils.CheckSemanticVersionString("0.1.0"))
	})

	t.Run("empty string", func(t *testing.T) {
		require.Error(t, versionutils.CheckSemanticVersionString(""))
	})

	t.Run("only major version", func(t *testing.T) {
		require.Error(t, versionutils.CheckSemanticVersionString("1"))
	})

	t.Run("only major and minor", func(t *testing.T) {
		require.Error(t, versionutils.CheckSemanticVersionString("1.2"))
	})

	t.Run("four parts", func(t *testing.T) {
		require.Error(t, versionutils.CheckSemanticVersionString("1.2.3.4"))
	})

	t.Run("with v prefix", func(t *testing.T) {
		require.NoError(t, versionutils.CheckSemanticVersionString("v1.2.3"))
	})

	t.Run("letters in version", func(t *testing.T) {
		require.Error(t, versionutils.CheckSemanticVersionString("a.b.c"))
	})

	t.Run("with pre-release suffix", func(t *testing.T) {
		require.Error(t, versionutils.CheckSemanticVersionString("1.2.3-beta"))
	})

	t.Run("with build metadata", func(t *testing.T) {
		require.Error(t, versionutils.CheckSemanticVersionString("1.2.3+build"))
	})

	t.Run("negative numbers", func(t *testing.T) {
		require.Error(t, versionutils.CheckSemanticVersionString("-1.2.3"))
	})

	t.Run("spaces around", func(t *testing.T) {
		require.Error(t, versionutils.CheckSemanticVersionString(" 1.2.3 "))
	})

	t.Run("date version is not a semantic version", func(t *testing.T) {
		require.Error(t, versionutils.CheckSemanticVersionString("20260725_143059"))
	})
}

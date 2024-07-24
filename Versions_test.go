package asciichgolangpublic

import (
	"testing"


	"github.com/stretchr/testify/assert"
)

func TestVersionsGetDateVersionString(t *testing.T) {
	tests := []struct {
		testmessage string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				dateVersion := Versions().MustGetNewDateVersionString()
				assert.Len(dateVersion, len("YYYYmmdd_HHMMSS"))

				assert.True(Versions().MustCheckDateVersionString(dateVersion))
			},
		)
	}
}

func TestVersionsGetSoftwareVersionEnvVarName(t *testing.T) {
	tests := []struct {
		testmessage string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues("SOFTWARE_VERSION", Versions().GetSoftwareVersionEnvVarName())
			},
		)
	}
}

func TestVersionsIsDateVersionString(t *testing.T) {
	tests := []struct {
		versionString           string
		expectedIsVersionString bool
	}{
		{"testcase", false},
		{"20231112_123456", true},
		{"20231112_12345", false},
		{"20231112_1234566", false},
		{"20231112_a23456", false},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedIsVersionString,
					Versions().IsDateVersionString(tt.versionString),
				)
			},
		)
	}
}

func TestVersionsIsVersionString(t *testing.T) {
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedIsVersionString,
					Versions().IsVersionString(tt.versionString),
				)
			},
		)
	}
}

func TestVersionsIsSemanticVersion(t *testing.T) {
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				var version Version = Versions().MustGetNewVersionByString(tt.versionString)

				assert.EqualValues(
					tt.expectedIsSemanticVersion,
					version.IsSemanticVersion(),
				)
			},
		)
	}
}

func TestVersionsGetLatestVersionFromSlice(t *testing.T) {
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				versions := Versions().MustGetVersionsFromStringSlice(tt.versionStrings)

				latestVersion := Versions().MustGetLatestVersionFromSlice(versions)

				expectedNewestVersion := MustGetVersionByString(tt.expectedNewest)

				assert.True(latestVersion.Equals(expectedNewestVersion))
			},
		)
	}
}

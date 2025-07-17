package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestDatesLayoutStringYYYYmmdd_HHMMSSS(t *testing.T) {

	require := require.New(t)

	require.EqualValues(
		"20060102_150405",
		Dates().LayoutStringYYYYmmdd_HHMMSS(),
	)
}

func TestDatesGetAsYYYYmmdd_HHMMSSString(t *testing.T) {
	tests := []struct {
		input    time.Time
		expected string
	}{
		{time.Date(2023, 11, 21, 14, 04, 24, 0, time.UTC), "20231121_140424"},
		{time.Date(2023, 11, 21, 2, 04, 24, 0, time.UTC), "20231121_020424"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				formated := Dates().MustFormatAsYYYYmmdd_HHMMSSString(&tt.input)

				require.EqualValues(tt.expected, formated)
			},
		)
	}
}

func TestDatesParseString_UTC(t *testing.T) {
	tests := []struct {
		input        string
		expectedDate time.Time
	}{
		{"Tue Nov 21 02:04:24 PM UTC 2023", time.Date(2023, 11, 21, 14, 04, 24, 0, time.UTC)},
		{"Tue Nov 21 02:04:24 AM UTC 2023", time.Date(2023, 11, 21, 2, 04, 24, 0, time.UTC)},
		{"20231121_142207", time.Date(2023, 11, 21, 14, 22, 07, 0, time.UTC)},
		{"20231121-142207", time.Date(2023, 11, 21, 14, 22, 07, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				parsed := Dates().MustParseString(tt.input)
				require.EqualValues(tt.expectedDate, parsed.UTC())
			},
		)
	}
}

func TestDatesParseString_CET(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		input        string
		expectedDate time.Time
	}{
		{"Tue Nov 21 02:04:24 PM CET 2023", time.Date(2023, 11, 21, 13, 04, 24, 0, time.UTC)},
		{"Tue Nov 21 02:04:24 AM CET 2023", time.Date(2023, 11, 21, 1, 04, 24, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				parsed := Dates().MustParseString(tt.input)
				require.EqualValues(tt.expectedDate, parsed.UTC())
			},
		)
	}
}

package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/testutils"
)

func TestDurationParserToSecondsInt64(t *testing.T) {
	tests := []struct {
		durationString   string
		expectedDuration int64
	}{
		{"1", 1},
		{"2", 2},
		{"3", 3},
		{"1s", 1},
		{"1second", 1},
		{"1seconds", 1},
		{"1 ", 1},
		{"1 s", 1},
		{"1 second", 1},
		{"1 seconds", 1},
		{" 1 ", 1},
		{" 1 s", 1},
		{" 1 second", 1},
		{" 1 seconds", 1},
		{" 1  ", 1},
		{" 1 s ", 1},
		{" 1 second ", 1},
		{" 1 seconds ", 1},
		{" 1\n", 1},
		{" 1 s\n", 1},
		{" 1 second\n", 1},
		{" 1 seconds\n", 1},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				parsedDuration := DurationParser().MustToSecondsInt64(tt.durationString)
				require.EqualValues(tt.expectedDuration, parsedDuration)
			},
		)
	}
}

func TestDurationParserToSecondsFloat64(t *testing.T) {
	tests := []struct {
		durationString   string
		expectedDuration float64
	}{
		{"1", 1.0},
		{"2", 2.0},
		{"3", 3.0},
		{"1s", 1.0},
		{"1second", 1.0},
		{"1seconds", 1.0},
		{"1 ", 1.0},
		{"1 s", 1.0},
		{"1 second", 1.0},
		{"1 seconds", 1.0},
		{" 1 ", 1.0},
		{" 1 s", 1.0},
		{" 1 second", 1.0},
		{" 1 seconds", 1.0},
		{" 1  ", 1.0},
		{" 1 s ", 1.0},
		{" 1 second ", 1.0},
		{" 1 seconds ", 1.0},
		{" 1\n", 1.0},
		{" 1 s\n", 1.0},
		{" 1 second\n", 1.0},
		{" 1 seconds\n", 1.0},
		{"0.5s", 0.5},
		{"0.5second", 0.5},
		{"0.5seconds", 0.5},
		{"1m", 60.0},
		{"1minute", 60.0},
		{"1minutes", 60.0},
		{"1 m", 60.0},
		{"1 minute", 60.0},
		{"1 minutes", 60.0},
		{"1h", 60.0 * 60.0},
		{"1hour", 60.0 * 60.0},
		{"1hours", 60.0 * 60.0},
		{"1 h", 60.0 * 60.0},
		{"1 hour", 60.0 * 60.0},
		{"1 hours", 60.0 * 60.0},
		{"1d", 60.0 * 60.0 * 24.0},
		{"1day", 60.0 * 60.0 * 24.0},
		{"1days", 60.0 * 60.0 * 24.0},
		{"1 d", 60.0 * 60.0 * 24.0},
		{"1 day", 60.0 * 60.0 * 24.0},
		{"1 days", 60.0 * 60.0 * 24.0},
		{"1w", 60.0 * 60.0 * 24.0 * 7.0},
		{"1week", 60.0 * 60.0 * 24.0 * 7.0},
		{"1weeks", 60.0 * 60.0 * 24.0 * 7.0},
		{"1 w", 60.0 * 60.0 * 24.0 * 7.0},
		{"1 week", 60.0 * 60.0 * 24.0 * 7.0},
		{"1 weeks", 60.0 * 60.0 * 24.0 * 7.0},
		{"1month", 60.0 * 60.0 * 24.0 * 30.0},
		{"1months", 60.0 * 60.0 * 24.0 * 30.0},
		{"1 month", 60.0 * 60.0 * 24.0 * 30.0},
		{"1 months", 60.0 * 60.0 * 24.0 * 30.0},
		{"1y", 60.0 * 60.0 * 24.0 * 364.0},
		{"1year", 60.0 * 60.0 * 24.0 * 364.0},
		{"1years", 60.0 * 60.0 * 24.0 * 364.0},
		{"1 y", 60.0 * 60.0 * 24.0 * 364.0},
		{"1 year", 60.0 * 60.0 * 24.0 * 364.0},
		{"1 years", 60.0 * 60.0 * 24.0 * 364.0},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				parsedDuration := DurationParser().MustToSecondsFloat64(tt.durationString)
				require.EqualValues(tt.expectedDuration, parsedDuration)
			},
		)
	}
}

func TestDurationParserToTimeDuration(t *testing.T) {
	tests := []struct {
		durationString   string
		expectedDuration time.Duration
	}{
		{"1", time.Second * 1},
		{"2", time.Second * 2},
		{"3", time.Second * 3},
		{"1s", time.Second * 1},
		{"1second", time.Second * 1},
		{"1seconds", time.Second * 1},
		{"1 ", time.Second * 1},
		{"1 s", time.Second * 1},
		{"1 second", time.Second * 1},
		{"1 seconds", time.Second * 1},
		{" 1 ", time.Second * 1},
		{" 1 s", time.Second * 1},
		{" 1 second", time.Second * 1},
		{" 1 seconds", time.Second * 1},
		{" 1  ", time.Second * 1},
		{" 1 s ", time.Second * 1},
		{" 1 second ", time.Second * 1},
		{" 1 seconds ", time.Second * 1},
		{" 1\n", time.Second * 1},
		{" 1 s\n", time.Second * 1},
		{" 1 second\n", time.Second * 1},
		{" 1 seconds\n", time.Second * 1},
		{"0.5s", time.Millisecond * 500},
		{"0.5second", time.Millisecond * 500},
		{"0.5seconds", time.Millisecond * 500},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				parsedDuration := DurationParser().MustToSecondsAsTimeDuration(tt.durationString)
				require.EqualValues(tt.expectedDuration, *parsedDuration)
			},
		)
	}
}

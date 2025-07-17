package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDurationFormatterToString(t *testing.T) {
	tests := []struct {
		duration         time.Duration
		expectedDuration string
	}{
		{time.Second * 0, "0s"},
		{time.Second * 1, "1s"},
		{time.Second * 2, "2s"},
		{time.Second * 12, "12s"},
		{time.Second * 60, "1m"},
		{time.Second * 90, "1m30s"},
		{time.Second * 91, "1m31s"},
		{time.Second * 60 * 10, "10m"},
		{time.Second * 60 * 59, "59m"},
		{time.Second * 60 * 60, "1h"},
		{time.Second * 60 * 61, "1h1m"},
		{time.Second*60*60 + time.Second*1, "1h1s"},
		{time.Second*60*62 + time.Second*1, "1h2m1s"},
		{time.Second * 60 * 60 * 23, "23h"},
		{time.Second * 60 * 60 * 24, "1d"},
		{time.Second * 60 * 60 * 24 * 29, "29d"},
		{time.Second * 60 * 60 * 24 * 30, "1months"},
		{time.Second * 60 * 60 * 24 * 30 * 3, "3months"},
		{time.Second * 60 * 60 * 24 * 364, "1year"},
	}

	for _, tt := range tests {
		t.Run(
			tt.expectedDuration,
			func(t *testing.T) {
				require := require.New(t)

				durationString := MustFormatDurationAsString(&tt.duration)
				require.EqualValues(tt.expectedDuration, durationString)
			},
		)
	}
}

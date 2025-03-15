package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
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
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				durationString := MustFormatDurationAsString(&tt.duration)
				require.EqualValues(tt.expectedDuration, durationString)
			},
		)
	}
}

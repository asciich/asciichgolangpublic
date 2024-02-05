package asciichgolangpublic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				durationString := DurationFormatter().MustToString(&tt.duration)
				assert.EqualValues(tt.expectedDuration, durationString)
			},
		)
	}
}

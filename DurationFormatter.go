package asciichgolangpublic

import (
	"fmt"
	"strings"
	"time"
)

type DurationFormatterService struct {
}

func DurationFormatter() (d *DurationFormatterService) {
	return NewDurationFormatterService()
}

func NewDurationFormatterService() (d *DurationFormatterService) {
	return new(DurationFormatterService)
}

func (d *DurationFormatterService) MustToString(duration *time.Duration) (durationString string) {
	durationString, err := d.ToString(duration)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return durationString
}

func (d *DurationFormatterService) ToString(duration *time.Duration) (durationString string, err error) {
	if duration == nil {
		return "", TracedError("duration is nil")
	}

	durationString = fmt.Sprintf("%v", *duration)

	if durationString != "0s" {
		durationString = strings.TrimSuffix(durationString, "0s")
	}

	return durationString, nil
}

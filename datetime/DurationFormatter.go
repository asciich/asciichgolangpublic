package datetime

import (
	"fmt"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func MustFormatDurationAsString(duration *time.Duration) (durationString string) {
	durationString, err := FormatDurationAsString(duration)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return durationString
}

func FormatDurationAsString(duration *time.Duration) (durationString string, err error) {
	if duration == nil {
		return "", tracederrors.TracedError("duration is nil")
	}

	durationString = fmt.Sprintf("%v", *duration)

	if durationString != "0s" {
		durationString = strings.TrimSuffix(durationString, "0s")
	}

	return durationString, nil
}

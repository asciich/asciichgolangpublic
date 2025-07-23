package datetime

import (
	"fmt"
	"time"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
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

	rest := int64(*duration)
	if rest == 0 {
		return "0s", nil
	}

	type timeUnits struct {
		Unit           string
		SecondsPerUnit int64
	}

	units := []timeUnits{
		{"year", int64(time.Second * 60 * 60 * 24 * 364)},
		{"months", int64(time.Second * 60 * 60 * 24 * 30)},
		{"d", int64(time.Second * 60 * 60 * 24)},
		{"h", int64(time.Second * 60 * 60)},
		{"m", int64(time.Second * 60)},
	}

	for _, u := range units {
		if rest >= u.SecondsPerUnit {
			v := rest / u.SecondsPerUnit
			rest = rest % u.SecondsPerUnit
			durationString += fmt.Sprintf("%d%s", v, u.Unit)
		}
	}

	if rest != 0 {
		durationString += fmt.Sprintf("%v", time.Duration(rest))
	}

	return durationString, nil
}

package datetime

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type DurationParserService struct{}

func DurationParser() (durationParser *DurationParserService) {
	return new(DurationParserService)
}

func NewDurationParserService() (d *DurationParserService) {
	return new(DurationParserService)
}

func (d *DurationParserService) MustToSecondsAsString(durationString string) (secondsString string) {
	secondsString, err := d.ToSecondsAsString(durationString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return secondsString
}

func (d *DurationParserService) MustToSecondsAsTimeDuration(durationString string) (duration *time.Duration) {
	duration, err := d.ToSecondsAsTimeDuration(durationString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return duration
}

func (d *DurationParserService) MustToSecondsFloat64(durationString string) (seconds float64) {
	seconds, err := d.ToSecondsFloat64(durationString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return seconds
}

func (d *DurationParserService) MustToSecondsInt64(durationString string) (seconds int64) {
	seconds, err := d.ToSecondsInt64(durationString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return seconds
}

func (d *DurationParserService) ToSecondsAsString(durationString string) (secondsString string, err error) {
	seconds, err := d.ToSecondsInt64(durationString)
	if err != nil {
		return "", err
	}

	secondsString = fmt.Sprintf("%v", seconds)
	return secondsString, nil
}

func (d *DurationParserService) ToSecondsAsTimeDuration(durationString string) (duration *time.Duration, err error) {
	secondsFloat, err := d.ToSecondsFloat64(durationString)
	if err != nil {
		return nil, err
	}

	durationToReturn := time.Duration(int(secondsFloat * float64(time.Second)))
	return &durationToReturn, nil
}

func (d *DurationParserService) ToSecondsFloat64(durationString string) (seconds float64, err error) {
	if len(durationString) <= 0 {
		return -1, tracederrors.TracedError("durationString is empty string")
	}

	unifiedDurationString := durationString
	unifiedDurationString = strings.TrimSpace(unifiedDurationString)
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "seconds", "s")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "second", "s")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "minutes", "m")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "minute", "m")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "hours", "h")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "hour", "h")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "days", "d")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "day", "d")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "weeks", "w")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "week", "w")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "months", "x")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "month", "x")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "years", "y")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, "year", "y")
	unifiedDurationString = strings.ReplaceAll(unifiedDurationString, " ", "")

	unitsAndSeconds := map[string]float64{
		"s": 1.0,
		"m": 60.0,
		"h": 60.0 * 60.0,
		"d": 60.0 * 60.0 * 24.0,
		"w": 60.0 * 60.0 * 24.0 * 7.0,
		"x": 60.0 * 60.0 * 24.0 * 30.0,
		"y": 60.0 * 60.0 * 24.0 * 364.0,
	}

	seconds = 0.0
	for k, v := range unitsAndSeconds {
		if !strings.HasSuffix(unifiedDurationString, k) {
			continue
		}

		unifiedDurationString = strings.TrimSuffix(unifiedDurationString, k)

		parsedValue, err := strconv.ParseFloat(unifiedDurationString, 64)
		if err != nil {
			return -1, tracederrors.TracedError(err.Error())
		}

		seconds = parsedValue * v
		return seconds, nil
	}

	seconds, err = strconv.ParseFloat(unifiedDurationString, 64)
	if err != nil {
		return -1, tracederrors.TracedError(err.Error())
	}

	return seconds, nil
}

func (d *DurationParserService) ToSecondsInt64(durationString string) (seconds int64, err error) {
	secondsFloat, err := d.ToSecondsFloat64(durationString)
	if err != nil {
		return -1, err
	}

	seconds = int64(secondsFloat)
	return seconds, nil
}

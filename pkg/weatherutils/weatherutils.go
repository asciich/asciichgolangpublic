package weatherutils

import (
	"context"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/weatherutils/openmeteoutils"
)

// GetCurrentTemperature returns the current temperature in degree celsius at the given coordinates.
//
// This is a convenience function delegating to [openmeteoutils.GetCurrentTemperature].
// The data is provided by open-meteo.com. It is free for non-commercial use, needs no API key
// and is licensed under CC BY 4.0. A credit message is logged after every successful request.
func GetCurrentTemperature(ctx context.Context, lat float64, lon float64) (float64, error) {
	return openmeteoutils.GetCurrentTemperature(ctx, lat, lon)
}

// GetHourlyTemperatures returns the temperature in degree celsius of every full hour
// from start to end (both inclusive, UTC) at the given coordinates.
//
// This is a convenience function delegating to [openmeteoutils.GetHourlyTemperatures].
//
// No root privileges are required.
func GetHourlyTemperatures(ctx context.Context, lat float64, lon float64, start time.Time, end time.Time) ([]openmeteoutils.HourlyTemperature, error) {
	return openmeteoutils.GetHourlyTemperatures(ctx, lat, lon, start, end)
}

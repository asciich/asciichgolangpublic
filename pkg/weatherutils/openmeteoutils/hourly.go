package openmeteoutils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// The time format of the hourly values when requested with timezone=UTC.
const hourlyTimeFormat = "2006-01-02T15:04"

// HourlyTemperature is the temperature of one full hour.
type HourlyTemperature struct {
	// The full hour in UTC.
	Time time.Time

	// Temperature in degree celsius.
	Temperature float64
}

type hourlyResponse struct {
	Hourly *struct {
		Time          []string   `json:"time"`
		Temperature2m []*float64 `json:"temperature_2m"`
	} `json:"hourly"`
}

// GetHourlyTemperatures requests the temperature of every full hour from start to end (both inclusive)
// at the given coordinates from open-meteo.com. The temperatures are in degree celsius.
//
// start and end are rounded down to the full hour and handled in UTC.
// The result is sorted by time. Hours without a value are skipped and logged.
// It is an error if there is no value at all.
//
// This is done with a single request. The forecast endpoint only returns data
// for roughly the last 92 days. For older data the Historical Weather API is needed.
//
// The values are modeled data for the given coordinates and not a measurement of a local sensor.
// The data is licensed under CC BY 4.0. Therefore a credit message is logged after every successful request.
//
// No root privileges are required.
func GetHourlyTemperatures(ctx context.Context, lat float64, lon float64, start time.Time, end time.Time) ([]HourlyTemperature, error) {
	err := checkCoordinates(lat, lon)
	if err != nil {
		return nil, err
	}

	start = start.UTC().Truncate(time.Hour)
	end = end.UTC().Truncate(time.Hour)

	if end.Before(start) {
		return nil, tracederrors.TracedErrorf("The end '%s' is before the start '%s'.", end.Format(hourlyTimeFormat), start.Format(hourlyTimeFormat))
	}

	logging.LogInfoByCtxf(
		ctx,
		"Get hourly temperatures at latitude '%v' and longitude '%v' from '%s' to '%s' (UTC) from open-meteo.com started.",
		lat, lon, start.Format(hourlyTimeFormat), end.Format(hourlyTimeFormat),
	)

	query := newForecastQuery(lat, lon)
	query.Set("hourly", "temperature_2m")
	query.Set("start_date", start.Format("2006-01-02"))
	query.Set("end_date", end.Format("2006-01-02"))

	body, err := requestForecast(ctx, query)
	if err != nil {
		return nil, err
	}

	ret, skipped, err := parseHourlyTemperatures(body, start, end)
	if err != nil {
		return nil, err
	}

	if skipped > 0 {
		logging.LogWarnByCtxf(ctx, "Open-meteo.com returned no value for '%d' hours. They are skipped.", skipped)
	}

	logging.LogInfoByCtxf(ctx, "Get hourly temperatures from open-meteo.com finished. Got '%d' hourly values.", len(ret))
	logging.LogInfoByCtxf(ctx, creditMessage)

	return ret, nil
}

// parseHourlyTemperatures extracts the hourly temperatures between start and end (inclusive)
// from the response body and returns them together with the number of skipped hours without a value.
func parseHourlyTemperatures(body []byte, start time.Time, end time.Time) ([]HourlyTemperature, int, error) {
	var response hourlyResponse

	err := json.Unmarshal(body, &response)
	if err != nil {
		return nil, 0, tracederrors.TracedErrorf("Unable to parse open-meteo.com response as JSON: %w", err)
	}

	if response.Hourly == nil {
		return nil, 0, tracederrors.TracedErrorf("The open-meteo.com response has no 'hourly' section. Response body: '%s'", string(body))
	}

	if len(response.Hourly.Time) != len(response.Hourly.Temperature2m) {
		return nil, 0, tracederrors.TracedErrorf(
			"The open-meteo.com response has '%d' times but '%d' temperatures.",
			len(response.Hourly.Time),
			len(response.Hourly.Temperature2m),
		)
	}

	ret := []HourlyTemperature{}
	skipped := 0

	for i, timeString := range response.Hourly.Time {
		hour, err := time.ParseInLocation(hourlyTimeFormat, timeString, time.UTC)
		if err != nil {
			return nil, 0, tracederrors.TracedErrorf("Unable to parse the time '%s' of the open-meteo.com response: %w", timeString, err)
		}

		// The request covers whole days. Only the wanted hours are returned.
		if hour.Before(start) || hour.After(end) {
			continue
		}

		if response.Hourly.Temperature2m[i] == nil {
			skipped++
			continue
		}

		ret = append(ret, HourlyTemperature{Time: hour, Temperature: *response.Hourly.Temperature2m[i]})
	}

	if len(ret) == 0 {
		return nil, 0, tracederrors.TracedErrorf("The open-meteo.com response has no temperature value between '%s' and '%s'.", start.Format(hourlyTimeFormat), end.Format(hourlyTimeFormat))
	}

	return ret, skipped, nil
}

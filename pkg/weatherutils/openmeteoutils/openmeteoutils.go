package openmeteoutils

import (
	"context"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	forecastApiUrl = "https://api.open-meteo.com/v1/forecast"
	requestTimeout = 30 * time.Second

	// Upper limit for the response body to read.
	maxResponseBodyBytes = 1 << 20

	// The data of open-meteo.com is licensed under CC BY 4.0 which requires to credit the source.
	creditMessage = "Weather data by Open-Meteo.com (https://open-meteo.com/), licensed under CC BY 4.0 (https://creativecommons.org/licenses/by/4.0/)."
)

// The relevant part of the open-meteo.com forecast response.
// Pointers are used to distinguish a missing value from a temperature of exactly 0 °C.
type forecastResponse struct {
	Current *struct {
		Temperature2m *float64 `json:"temperature_2m"`
	} `json:"current"`
}

// GetCurrentTemperature requests the current temperature at the given coordinates from open-meteo.com
// and returns it in degree celsius.
//
// The open-meteo.com API is free for non-commercial use and needs no API key.
// The data is licensed under CC BY 4.0. Therefore a credit message is logged
// after every successful request.
//
// The returned value is modeled data for the given coordinates and not a measurement of a local sensor.
//
// No root privileges are required.
func GetCurrentTemperature(ctx context.Context, lat float64, lon float64) (float64, error) {
	err := checkCoordinates(lat, lon)
	if err != nil {
		return 0, err
	}

	logging.LogInfoByCtxf(ctx, "Get current temperature at latitude '%v' and longitude '%v' from open-meteo.com started.", lat, lon)

	query := newForecastQuery(lat, lon)
	query.Set("current", "temperature_2m")

	body, err := requestForecast(ctx, query)
	if err != nil {
		return 0, err
	}

	ret, err := parseCurrentTemperature(body)
	if err != nil {
		return 0, err
	}

	logging.LogInfoByCtxf(ctx, "Get current temperature at latitude '%v' and longitude '%v' from open-meteo.com finished. Temperature is '%.1f' °C.", lat, lon, ret)
	logging.LogInfoByCtxf(ctx, creditMessage)

	return ret, nil
}

func checkCoordinates(lat float64, lon float64) error {
	if math.IsNaN(lat) || lat < -90 || lat > 90 {
		return tracederrors.TracedErrorf("Invalid latitude '%v'. The latitude must be between -90 and 90.", lat)
	}

	if math.IsNaN(lon) || lon < -180 || lon > 180 {
		return tracederrors.TracedErrorf("Invalid longitude '%v'. The longitude must be between -180 and 180.", lon)
	}

	return nil
}

// newForecastQuery returns the query parameters shared by all requests.
func newForecastQuery(lat float64, lon float64) url.Values {
	ret := url.Values{}
	ret.Set("latitude", strconv.FormatFloat(lat, 'f', -1, 64))
	ret.Set("longitude", strconv.FormatFloat(lon, 'f', -1, 64))
	ret.Set("temperature_unit", "celsius")
	ret.Set("timezone", "UTC")

	return ret
}

// requestForecast performs the request and returns the response body of a successful (HTTP 200) answer.
func requestForecast(ctx context.Context, query url.Values) ([]byte, error) {
	requestUrl := forecastApiUrl + "?" + query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to create request for '%s': %w", requestUrl, err)
	}

	client := &http.Client{Timeout: requestTimeout}

	resp, err := client.Do(req)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Request to '%s' failed: %w", forecastApiUrl, err)
	}
	defer resp.Body.Close()

	ret, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyBytes))
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to read response body of '%s': %w", forecastApiUrl, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, tracederrors.TracedErrorf(
			"Request to '%s' returned unexpected status code '%d'. Response body: '%s'",
			forecastApiUrl,
			resp.StatusCode,
			string(ret),
		)
	}

	return ret, nil
}

// parseCurrentTemperature extracts the current temperature in degree celsius from the open-meteo.com forecast response body.
func parseCurrentTemperature(body []byte) (float64, error) {
	var response forecastResponse

	err := json.Unmarshal(body, &response)
	if err != nil {
		return 0, tracederrors.TracedErrorf("Unable to parse open-meteo.com response as JSON: %w", err)
	}

	if response.Current == nil {
		return 0, tracederrors.TracedErrorf("The open-meteo.com response has no 'current' section. Response body: '%s'", string(body))
	}

	if response.Current.Temperature2m == nil {
		return 0, tracederrors.TracedErrorf("The open-meteo.com response has no 'current.temperature_2m' value. Response body: '%s'", string(body))
	}

	ret := *response.Current.Temperature2m

	return ret, nil
}

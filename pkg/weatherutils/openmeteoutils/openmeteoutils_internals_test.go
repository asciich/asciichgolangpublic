package openmeteoutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_parseCurrentTemperature(t *testing.T) {
	t.Run("positive temperature", func(t *testing.T) {
		temperature, err := parseCurrentTemperature([]byte(`{"current":{"time":"2026-11-21T19:45","interval":900,"temperature_2m":6.3}}`))
		require.NoError(t, err)
		require.InDelta(t, 6.3, temperature, 0.0001)
	})

	t.Run("negative temperature", func(t *testing.T) {
		temperature, err := parseCurrentTemperature([]byte(`{"current":{"temperature_2m":-2.5}}`))
		require.NoError(t, err)
		require.InDelta(t, -2.5, temperature, 0.0001)
	})

	t.Run("exactly zero is a valid temperature", func(t *testing.T) {
		temperature, err := parseCurrentTemperature([]byte(`{"current":{"temperature_2m":0}}`))
		require.NoError(t, err)
		require.InDelta(t, 0.0, temperature, 0.0001)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := parseCurrentTemperature([]byte(`not json`))
		require.Error(t, err)
	})

	t.Run("missing current section", func(t *testing.T) {
		_, err := parseCurrentTemperature([]byte(`{"error":true,"reason":"something went wrong"}`))
		require.Error(t, err)
	})

	t.Run("missing temperature value", func(t *testing.T) {
		_, err := parseCurrentTemperature([]byte(`{"current":{"time":"2026-11-21T19:45"}}`))
		require.Error(t, err)
	})
}

func Test_parseHourlyTemperatures(t *testing.T) {
	start := time.Date(2026, 11, 18, 1, 0, 0, 0, time.UTC)
	end := time.Date(2026, 11, 18, 3, 0, 0, 0, time.UTC)

	t.Run("only requested hours, null skipped", func(t *testing.T) {
		body := []byte(`{"hourly":{"time":["2026-11-18T00:00","2026-11-18T01:00","2026-11-18T02:00","2026-11-18T03:00","2026-11-18T04:00"],"temperature_2m":[1.0,2.5,null,0,9.0]}}`)

		temperatures, skipped, err := parseHourlyTemperatures(body, start, end)
		require.NoError(t, err)
		require.Equal(t, 1, skipped)
		require.Len(t, temperatures, 2)
		require.Equal(t, start, temperatures[0].Time)
		require.InDelta(t, 2.5, temperatures[0].Temperature, 0.0001)
		require.InDelta(t, 0.0, temperatures[1].Temperature, 0.0001)
	})

	t.Run("no values at all", func(t *testing.T) {
		body := []byte(`{"hourly":{"time":["2026-11-18T01:00"],"temperature_2m":[null]}}`)

		_, _, err := parseHourlyTemperatures(body, start, end)
		require.Error(t, err)
	})

	t.Run("length mismatch", func(t *testing.T) {
		body := []byte(`{"hourly":{"time":["2026-11-18T01:00","2026-11-18T02:00"],"temperature_2m":[1.0]}}`)

		_, _, err := parseHourlyTemperatures(body, start, end)
		require.Error(t, err)
	})

	t.Run("missing hourly section", func(t *testing.T) {
		_, _, err := parseHourlyTemperatures([]byte(`{"error":true,"reason":"x"}`), start, end)
		require.Error(t, err)
	})
}

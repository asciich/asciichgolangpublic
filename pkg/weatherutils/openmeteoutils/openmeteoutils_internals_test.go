package openmeteoutils

import (
	"testing"

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

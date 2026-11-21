package openmeteoutils_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/weatherutils/openmeteoutils"
)

func Test_GetCurrentTemperature(t *testing.T) {
	t.Run("latitude NaN", func(t *testing.T) {
		_, err := openmeteoutils.GetCurrentTemperature(context.Background(), math.NaN(), 7.9)
		require.Error(t, err)
	})

	t.Run("latitude too small", func(t *testing.T) {
		_, err := openmeteoutils.GetCurrentTemperature(context.Background(), -90.1, 7.9)
		require.Error(t, err)
	})

	t.Run("latitude too large", func(t *testing.T) {
		_, err := openmeteoutils.GetCurrentTemperature(context.Background(), 90.1, 7.9)
		require.Error(t, err)
	})

	t.Run("longitude NaN", func(t *testing.T) {
		_, err := openmeteoutils.GetCurrentTemperature(context.Background(), 47.3, math.NaN())
		require.Error(t, err)
	})

	t.Run("longitude too small", func(t *testing.T) {
		_, err := openmeteoutils.GetCurrentTemperature(context.Background(), 47.3, -180.1)
		require.Error(t, err)
	})

	t.Run("longitude too large", func(t *testing.T) {
		_, err := openmeteoutils.GetCurrentTemperature(context.Background(), 47.3, 180.1)
		require.Error(t, err)
	})

	t.Run("cancelled context returns error", func(t *testing.T) {
		// A cancelled context aborts before any data is sent, so no real request is performed.
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := openmeteoutils.GetCurrentTemperature(ctx, 47.36583, 7.903761)
		require.Error(t, err)
	})

	// This is the only subtest performing a real request to open-meteo.com to avoid rate limiting.
	t.Run("real request returns plausible temperature", func(t *testing.T) {
		temperature, err := openmeteoutils.GetCurrentTemperature(context.Background(), 47.36583, 7.903761)
		require.NoError(t, err)

		// Plausibility check only: The real value changes all the time.
		require.Greater(t, temperature, -40.0)
		require.Less(t, temperature, 50.0)
	})

	t.Run("hourly: end before start", func(t *testing.T) {
		now := time.Now()
		_, err := openmeteoutils.GetHourlyTemperatures(context.Background(), 47.3, 7.9, now, now.Add(-2*time.Hour))
		require.Error(t, err)
	})
}

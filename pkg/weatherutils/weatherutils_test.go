package weatherutils_test

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/weatherutils"
)

func Test_GetCurrentTemperature(t *testing.T) {
	t.Run("invalid latitude", func(t *testing.T) {
		_, err := weatherutils.GetCurrentTemperature(context.Background(), 90.1, 7.9)
		require.Error(t, err)
	})

	t.Run("invalid longitude", func(t *testing.T) {
		_, err := weatherutils.GetCurrentTemperature(context.Background(), 47.3, math.NaN())
		require.Error(t, err)
	})

	// This is the only subtest performing a real request to open-meteo.com to avoid rate limiting.
	t.Run("real request returns plausible temperature", func(t *testing.T) {
		temperature, err := weatherutils.GetCurrentTemperature(context.Background(), 47.36583, 7.903761)
		require.NoError(t, err)

		// Plausibility check only: The real value changes all the time.
		require.Greater(t, temperature, -40.0)
		require.Less(t, temperature, 50.0)
	})
}

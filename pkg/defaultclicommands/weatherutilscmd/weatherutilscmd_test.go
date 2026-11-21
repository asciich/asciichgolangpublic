package weatherutilscmd_test

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/weatherutilscmd"
)

func Test_NewWeatherUtilsCmd(t *testing.T) {
	t.Run("root command is named weather", func(t *testing.T) {
		cmd := weatherutilscmd.NewWeatherUtilsCmd()
		require.Equal(t, "weather", cmd.Use)
	})

	t.Run("get-current-temperature is registered", func(t *testing.T) {
		cmd := weatherutilscmd.NewWeatherUtilsCmd()

		sub, _, err := cmd.Find([]string{"get-current-temperature"})
		require.NoError(t, err)
		require.Equal(t, "get-current-temperature", sub.Name())
	})
}

func Test_GetCurrentTemperatureCmd(t *testing.T) {
	t.Run("flags lat and lon exist", func(t *testing.T) {
		cmd := weatherutilscmd.NewGetCurrentTemperatureCmd()
		require.NotNil(t, cmd.Flags().Lookup("lat"))
		require.NotNil(t, cmd.Flags().Lookup("lon"))
	})

	t.Run("missing lat and lon", func(t *testing.T) {
		cmd := weatherutilscmd.NewGetCurrentTemperatureCmd()
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})
		cmd.SetArgs([]string{})

		require.Error(t, cmd.Execute())
	})

	t.Run("missing lon", func(t *testing.T) {
		cmd := weatherutilscmd.NewGetCurrentTemperatureCmd()
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})
		cmd.SetArgs([]string{"--lat", "47.36583"})

		require.Error(t, cmd.Execute())
	})

	t.Run("missing lat", func(t *testing.T) {
		cmd := weatherutilscmd.NewGetCurrentTemperatureCmd()
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})
		cmd.SetArgs([]string{"--lon", "7.903761"})

		require.Error(t, cmd.Execute())
	})

	t.Run("lat is not a number", func(t *testing.T) {
		cmd := weatherutilscmd.NewGetCurrentTemperatureCmd()
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})
		cmd.SetArgs([]string{"--lat", "abc", "--lon", "7.903761"})

		require.Error(t, cmd.Execute())
	})

	// This is the only subtest performing a real request to open-meteo.com to avoid rate limiting.
	t.Run("real request prints plausible temperature", func(t *testing.T) {
		cmd := weatherutilscmd.NewGetCurrentTemperatureCmd()

		stdout := &bytes.Buffer{}
		cmd.SetOut(stdout)
		cmd.SetErr(&bytes.Buffer{})
		cmd.SetArgs([]string{"--lat", "47.36583", "--lon", "7.903761"})

		require.NoError(t, cmd.Execute())

		temperature, err := strconv.ParseFloat(strings.TrimSpace(stdout.String()), 64)
		require.NoError(t, err)

		// Plausibility check only: The real value changes all the time.
		require.Greater(t, temperature, -40.0)
		require.Less(t, temperature, 50.0)
	})
}

package environmentvariables_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/environmentvariables"
)

func Test_GetEnvValueAsStringOrDefault(t *testing.T) {
	ctx := context.Background()

	t.Run("empty envName", func(t *testing.T) {
		_, err := environmentvariables.GetEnvValueAsStringOrDefault(ctx, "", "default")
		require.Error(t, err)
	})

	t.Run("env not set returns default", func(t *testing.T) {
		envName := "TEST_GET_ENV_VALUE_AS_STRING_OR_DEFAULT_NOT_SET"
		os.Unsetenv(envName)

		got, err := environmentvariables.GetEnvValueAsStringOrDefault(ctx, envName, "myDefault")
		require.NoError(t, err)
		require.EqualValues(t, "myDefault", got)
	})

	t.Run("env set to empty returns default", func(t *testing.T) {
		envName := "TEST_GET_ENV_VALUE_AS_STRING_OR_DEFAULT_EMPTY"
		os.Setenv(envName, "")
		defer os.Unsetenv(envName)

		got, err := environmentvariables.GetEnvValueAsStringOrDefault(ctx, envName, "myDefault")
		require.NoError(t, err)
		require.EqualValues(t, "myDefault", got)
	})

	t.Run("env set returns value", func(t *testing.T) {
		envName := "TEST_GET_ENV_VALUE_AS_STRING_OR_DEFAULT_SET"
		os.Setenv(envName, "actualValue")
		defer os.Unsetenv(envName)

		got, err := environmentvariables.GetEnvValueAsStringOrDefault(ctx, envName, "myDefault")
		require.NoError(t, err)
		require.EqualValues(t, "actualValue", got)
	})

	t.Run("empty default value", func(t *testing.T) {
		envName := "TEST_GET_ENV_VALUE_AS_STRING_OR_DEFAULT_EMPTY_DEFAULT"
		os.Unsetenv(envName)

		got, err := environmentvariables.GetEnvValueAsStringOrDefault(ctx, envName, "")
		require.NoError(t, err)
		require.EqualValues(t, "", got)
	})
}

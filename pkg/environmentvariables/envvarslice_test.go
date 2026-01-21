package environmentvariables_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/environmentvariables"
)

func Test_SetEnvVarInStringSlice(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		got, err := environmentvariables.SetEnvVarInStringSlice(nil, "MY_VAL", "0")
		require.NoError(t, err)
		require.EqualValues(t, []string{"MY_VAL=0"}, got)
	})

	t.Run("empty", func(t *testing.T) {
		got, err := environmentvariables.SetEnvVarInStringSlice([]string{}, "MY_VAL", "0")
		require.NoError(t, err)
		require.EqualValues(t, []string{"MY_VAL=0"}, got)
	})

	t.Run("Add", func(t *testing.T) {
		got, err := environmentvariables.SetEnvVarInStringSlice([]string{"A=123"}, "MY_VAL", "0")
		require.NoError(t, err)
		require.EqualValues(t, []string{"A=123", "MY_VAL=0"}, got)
	})

	t.Run("Overwrite", func(t *testing.T) {
		got, err := environmentvariables.SetEnvVarInStringSlice([]string{"A=123", "MY_VAL=1"}, "MY_VAL", "0")
		require.NoError(t, err)
		require.EqualValues(t, []string{"A=123", "MY_VAL=0"}, got)
	})
}

func Test_SetEnvVarsInStringSlice(t *testing.T) {
	t.Run("nil nil", func(t *testing.T) {
		got, err := environmentvariables.SetEnvVarsInStringSlice(nil, nil)
		require.NoError(t, err)
		require.EqualValues(t, []string{}, got)
	})

	t.Run("nil", func(t *testing.T) {
		got, err := environmentvariables.SetEnvVarsInStringSlice(nil, map[string]string{"MY_VAL": "0"})
		require.NoError(t, err)
		require.EqualValues(t, []string{"MY_VAL=0"}, got)
	})

	t.Run("empty", func(t *testing.T) {
		got, err := environmentvariables.SetEnvVarsInStringSlice([]string{}, map[string]string{"MY_VAL": "0"})
		require.NoError(t, err)
		require.EqualValues(t, []string{"MY_VAL=0"}, got)
	})

	t.Run("Add", func(t *testing.T) {
		got, err := environmentvariables.SetEnvVarsInStringSlice([]string{"A=123"}, map[string]string{"MY_VAL": "0"})
		require.NoError(t, err)
		require.EqualValues(t, []string{"A=123", "MY_VAL=0"}, got)
	})

	t.Run("Overwrite", func(t *testing.T) {
		got, err := environmentvariables.SetEnvVarsInStringSlice([]string{"A=123", "MY_VAL=1"}, map[string]string{"MY_VAL": "0"})
		require.NoError(t, err)
		require.EqualValues(t, []string{"A=123", "MY_VAL=0"}, got)
	})
}

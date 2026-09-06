package parameteroptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func TestRunCommandOptions_GetCommand(t *testing.T) {
	// This test ensures the GetCommand returns:
	// - Only the original Command if set
	// - An error if unset or empty
	t.Run("unset", func(t *testing.T) {
		options := &parameteroptions.RunCommandOptions{}
		command, err := options.GetCommand()
		require.Error(t, err)
		require.Nil(t, command)
	})

	t.Run("empty", func(t *testing.T) {
		options := &parameteroptions.RunCommandOptions{
			Command: []string{},
		}
		command, err := options.GetCommand()
		require.Error(t, err)
		require.Nil(t, command)
	})

	tests := []struct {
		name     string
		options  *parameteroptions.RunCommandOptions
		expected []string
	}{
		{"single command", &parameteroptions.RunCommandOptions{Command: []string{"echo"}}, []string{"echo"}},
		{"hello world", &parameteroptions.RunCommandOptions{Command: []string{"echo", "hello", "world"}}, []string{"echo", "hello", "world"}},
		{"timeout", &parameteroptions.RunCommandOptions{TimeoutString: "1m", Command: []string{"echo", "hello", "world"}}, []string{"echo", "hello", "world"}},
		{"sudo timeout", &parameteroptions.RunCommandOptions{RunAsRoot: true, TimeoutString: "1m", Command: []string{"echo", "hello", "world"}}, []string{"echo", "hello", "world"}},
		{"timeout another user", &parameteroptions.RunCommandOptions{RunAsUser: "testuser", TimeoutString: "1m", Command: []string{"echo", "hello", "world"}}, []string{"echo", "hello", "world"}},
		{"sudo timeout another user", &parameteroptions.RunCommandOptions{UseSudoToRunAsUser: true, RunAsUser: "testuser", TimeoutString: "1m", Command: []string{"echo", "hello", "world"}}, []string{"echo", "hello", "world"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, err := tt.options.GetCommand()
			require.NoError(t, err)
			require.EqualValues(t, tt.expected, command)
		})
	}
}

func TestRunCommandOptions_GetFullCommand(t *testing.T) {
	// This test ensures the GetFullCommand returns:
	// - Only the original Command and additionally all prefix commands like 'sudo', 'timeout' when set in the options
	// - An error if unset or empty
	t.Run("unset", func(t *testing.T) {
		options := &parameteroptions.RunCommandOptions{}
		command, err := options.GetCommand()
		require.Error(t, err)
		require.Nil(t, command)
	})

	t.Run("empty", func(t *testing.T) {
		options := &parameteroptions.RunCommandOptions{
			Command: []string{},
		}
		command, err := options.GetCommand()
		require.Error(t, err)
		require.Nil(t, command)
	})

	tests := []struct {
		name     string
		options  *parameteroptions.RunCommandOptions
		expected []string
	}{
		{"single command", &parameteroptions.RunCommandOptions{Command: []string{"echo"}}, []string{"echo"}},
		{"hello world", &parameteroptions.RunCommandOptions{Command: []string{"echo", "hello", "world"}}, []string{"echo", "hello", "world"}},
		{"timeout", &parameteroptions.RunCommandOptions{TimeoutString: "1m", Command: []string{"echo", "hello", "world"}}, []string{"timeout", "60", "echo", "hello", "world"}},
		{"sudo timeout", &parameteroptions.RunCommandOptions{RunAsRoot: true, TimeoutString: "1m", Command: []string{"echo", "hello", "world"}}, []string{"timeout", "60", "sudo", "echo", "hello", "world"}},
		{"timeout another user", &parameteroptions.RunCommandOptions{RunAsUser: "testuser", TimeoutString: "1m", Command: []string{"echo", "hello", "world"}}, []string{"timeout", "60", "su", "testuser", "-c", "echo hello world"}},
		{"sudo timeout another user", &parameteroptions.RunCommandOptions{UseSudoToRunAsUser: true, RunAsUser: "testuser", TimeoutString: "1m", Command: []string{"echo", "hello", "world"}}, []string{"timeout", "60", "sudo", "su", "testuser", "-c", "echo hello world"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, err := tt.options.GetFullCommand()
			require.NoError(t, err)
			require.EqualValues(t, tt.expected, command)
		})
	}
}

func TestRunCommandOptions_GetDeepCopy(t *testing.T) {
	t.Run("nil fields", func(t *testing.T) {
		original := parameteroptions.NewRunCommandOptions()
		copy := original.GetDeepCopy()

		require.NotNil(t, copy)

		// Modify copy
		copy.RunAsRoot = true
		copy.TimeoutString = "10s"
		copy.StdinString = "input"
		copy.AllowAllExitCodes = true
		copy.RemoveLastLineIfEmpty = true
		copy.RunAsUser = "testuser"
		copy.UseSudoToRunAsUser = true

		// Original should be unchanged
		require.EqualValues(t, false, original.RunAsRoot)
		require.EqualValues(t, "", original.TimeoutString)
		require.EqualValues(t, "", original.StdinString)
		require.EqualValues(t, false, original.AllowAllExitCodes)
		require.EqualValues(t, false, original.RemoveLastLineIfEmpty)
		require.EqualValues(t, "", original.RunAsUser)
		require.EqualValues(t, false, original.UseSudoToRunAsUser)
	})

	t.Run("with Command slice", func(t *testing.T) {
		original := parameteroptions.NewRunCommandOptions()
		err := original.SetCommand([]string{"echo", "hello"})
		require.NoError(t, err)

		copy := original.GetDeepCopy()

		// Get commands
		origCmd, err := original.GetCommand()
		require.NoError(t, err)
		copyCmd, err := copy.GetCommand()
		require.NoError(t, err)

		require.EqualValues(t, origCmd, copyCmd)

		// Modify copy's command slice
		copy.Command[0] = "modified"
		copy.Command = append(copy.Command, "extra")

		// Original should be unchanged (deep copy of slice)
		origCmd, err = original.GetCommand()
		require.NoError(t, err)
		require.EqualValues(t, []string{"echo", "hello"}, origCmd)

		copyCmd, err = copy.GetCommand()
		require.NoError(t, err)
		require.EqualValues(t, []string{"modified", "hello", "extra"}, copyCmd)
	})

	t.Run("with AdditionalEnvVars map", func(t *testing.T) {
		original := parameteroptions.NewRunCommandOptions()
		original.AdditionalEnvVars = map[string]string{
			"KEY1": "value1",
			"KEY2": "value2",
		}

		copy := original.GetDeepCopy()

		require.EqualValues(t, original.AdditionalEnvVars, copy.AdditionalEnvVars)

		// Modify copy's map
		copy.AdditionalEnvVars["KEY1"] = "modified"
		copy.AdditionalEnvVars["KEY3"] = "value3"

		// Original should be unchanged (deep copy of map)
		require.EqualValues(t, map[string]string{
			"KEY1": "value1",
			"KEY2": "value2",
		}, original.AdditionalEnvVars)
		require.EqualValues(t, map[string]string{
			"KEY1": "modified",
			"KEY2": "value2",
			"KEY3": "value3",
		}, copy.AdditionalEnvVars)
	})

	t.Run("with all fields set", func(t *testing.T) {
		original := parameteroptions.NewRunCommandOptions()
		err := original.SetCommand([]string{"ls", "-la"})
		require.NoError(t, err)
		original.TimeoutString = "30s"
		original.StdinString = "input data"
		original.AllowAllExitCodes = true
		original.RemoveLastLineIfEmpty = true
		original.RunAsUser = "testuser"
		original.UseSudoToRunAsUser = true
		original.AdditionalEnvVars = map[string]string{
			"ENV_VAR": "value",
		}

		copy := original.GetDeepCopy()

		// Verify all fields are copied
		require.EqualValues(t, original.TimeoutString, copy.TimeoutString)
		require.EqualValues(t, original.StdinString, copy.StdinString)
		require.EqualValues(t, original.AllowAllExitCodes, copy.AllowAllExitCodes)
		require.EqualValues(t, original.RemoveLastLineIfEmpty, copy.RemoveLastLineIfEmpty)
		require.EqualValues(t, original.RunAsUser, copy.RunAsUser)
		require.EqualValues(t, original.UseSudoToRunAsUser, copy.UseSudoToRunAsUser)
		require.EqualValues(t, original.AdditionalEnvVars, copy.AdditionalEnvVars)

		// Modify copy
		copy.TimeoutString = "60s"
		copy.AdditionalEnvVars["ENV_VAR"] = "modified"

		// Original should be unchanged
		require.EqualValues(t, "30s", original.TimeoutString)
		require.EqualValues(t, map[string]string{"ENV_VAR": "value"}, original.AdditionalEnvVars)
	})
}

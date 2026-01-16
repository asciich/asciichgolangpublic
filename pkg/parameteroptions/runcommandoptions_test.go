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
		{"sudo timeout another user", &parameteroptions.RunCommandOptions{UseSudoToRunAsUser: true , RunAsUser:  "testuser", TimeoutString: "1m", Command: []string{"echo", "hello", "world"}}, []string{"echo", "hello", "world"}},
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
		{"sudo timeout another user", &parameteroptions.RunCommandOptions{UseSudoToRunAsUser: true , RunAsUser:  "testuser", TimeoutString: "1m", Command: []string{"echo", "hello", "world"}}, []string{"timeout", "60", "sudo", "su", "testuser", "-c", "echo hello world"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, err := tt.options.GetFullCommand()
			require.NoError(t, err)
			require.EqualValues(t, tt.expected, command)
		})
	}
}
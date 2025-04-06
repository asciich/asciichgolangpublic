package ansibleutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

func ctx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_ParseListHostsOutput(t *testing.T) {
	t.Run("unknown output", func(t *testing.T) {
		parsed, err := ParseCliOutput(ctx(), "x")
		require.Nil(t, parsed)
		require.ErrorIs(t, err, ErrUnknwnAnsibleCliOutput)
	})

	t.Run("single host", func(t *testing.T) {
		parsed, err := ParseCliOutput(ctx(), "  hosts (1):\n    one.example.net\n")
		require.NoError(t, err)
		require.EqualValues(
			t,
			[]string{"one.example.net"},
			parsed.inventory.MustListHostNames(),
		)
	})

	t.Run("two hosts", func(t *testing.T) {
		parsed, err := ParseCliOutput(ctx(), "  hosts (2):\n    one.example.net\n    two.example.net")
		require.NoError(t, err)
		require.EqualValues(
			t,
			[]string{"one.example.net", "two.example.net"},
			parsed.inventory.MustListHostNames(),
		)
	})
}

func Test_isListHostsOutput(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		require.EqualValues(t, false, isListHostsOutput(""))
	})

	t.Run("hosts in second line", func(t *testing.T) {
		require.EqualValues(t, false, isListHostsOutput("text\nhosts"))
	})

	t.Run("only hosts", func(t *testing.T) {
		require.EqualValues(t, false, isListHostsOutput("hosts"))
		require.EqualValues(t, false, isListHostsOutput(" hosts"))
		require.EqualValues(t, false, isListHostsOutput("  hosts"))
	})

	t.Run("invalid hosts", func(t *testing.T) {
		require.EqualValues(t, false, isListHostsOutput("hosts listed"))
	})

	t.Run("empty brackets", func(t *testing.T) {
		require.EqualValues(t, false, isListHostsOutput("hosts ()"))
	})

	t.Run("no tailing double point", func(t *testing.T) {
		require.EqualValues(t, false, isListHostsOutput("hosts (1)"))
	})

	t.Run("one host", func(t *testing.T) {
		require.EqualValues(t, true, isListHostsOutput(`  hosts (1):`))
	})
}

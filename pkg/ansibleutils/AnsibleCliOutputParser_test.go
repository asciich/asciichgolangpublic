package ansibleutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/ansibleutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
)

func ctx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_ParseListHostsOutput(t *testing.T) {
	t.Run("unknown output", func(t *testing.T) {
		parsed, err := ansibleutils.ParseCliOutput(ctx(), "x")
		require.Nil(t, parsed)
		require.ErrorIs(t, err, ansibleutils.ErrUnknwnAnsibleCliOutput)
	})

	t.Run("single host", func(t *testing.T) {
		parsed, err := ansibleutils.ParseCliOutput(ctx(), "  hosts (1):\n    one.example.net\n")
		require.NoError(t, err)
		require.EqualValues(
			t,
			[]string{"one.example.net"},
			mustutils.Must(parsed.ListHostNames()),
		)
	})

	t.Run("two hosts", func(t *testing.T) {
		parsed, err := ansibleutils.ParseCliOutput(ctx(), "  hosts (2):\n    one.example.net\n    two.example.net")
		require.NoError(t, err)
		require.EqualValues(
			t,
			[]string{"one.example.net", "two.example.net"},
			mustutils.Must(parsed.ListHostNames()),
		)
	})
}

func Test_isListHostsOutput(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		require.EqualValues(t, false, ansibleutils.IsListHostsOutput(""))
	})

	t.Run("hosts in second line", func(t *testing.T) {
		require.EqualValues(t, false, ansibleutils.IsListHostsOutput("text\nhosts"))
	})

	t.Run("only hosts", func(t *testing.T) {
		require.EqualValues(t, false, ansibleutils.IsListHostsOutput("hosts"))
		require.EqualValues(t, false, ansibleutils.IsListHostsOutput(" hosts"))
		require.EqualValues(t, false, ansibleutils.IsListHostsOutput("  hosts"))
	})

	t.Run("invalid hosts", func(t *testing.T) {
		require.EqualValues(t, false, ansibleutils.IsListHostsOutput("hosts listed"))
	})

	t.Run("empty brackets", func(t *testing.T) {
		require.EqualValues(t, false, ansibleutils.IsListHostsOutput("hosts ()"))
	})

	t.Run("no tailing double point", func(t *testing.T) {
		require.EqualValues(t, false, ansibleutils.IsListHostsOutput("hosts (1)"))
	})

	t.Run("one host", func(t *testing.T) {
		require.EqualValues(t, true, ansibleutils.IsListHostsOutput(`  hosts (1):`))
	})
}

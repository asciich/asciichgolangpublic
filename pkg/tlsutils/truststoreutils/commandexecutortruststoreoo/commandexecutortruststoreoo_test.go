package commandexecutortruststoreoo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/truststoreutils/commandexecutortruststoreoo"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/truststoreutils/truststoreinterfaces"
)

func TestNewCommandExecutorTrustStore(t *testing.T) {
	t.Run("nil command executor must return error", func(t *testing.T) {
		trustStore, err := commandexecutortruststoreoo.NewCommandExecutorTrustStore(nil, false)
		require.Error(t, err)
		require.Nil(t, trustStore)
		require.Contains(t, err.Error(), "commandExecutor")
	})
}

func TestCommandExecutorTrustStore_InterfaceCompliance(t *testing.T) {
	// Compile-time check that CommandExecutorTrustStore implements TrustStore interface
	var _ truststoreinterfaces.TrustStore = (*commandexecutortruststoreoo.CommandExecutorTrustStore)(nil)
}

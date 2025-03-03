package ansibleutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Constructor(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		output := NewAnsibleCliOutput()
		require.EqualValues(t, "in memory ansible cli output", output.Name())
	})

	t.Run("number of hosts", func(t *testing.T) {
		output := NewAnsibleCliOutput()
		require.EqualValues(t, 0, output.GetNumberOfHosts(ctx()))
	})
}

func Test_CreateInventory(t *testing.T) {
	output := NewAnsibleCliOutput()
	inventory := output.CreateInventory()
	require.EqualValues(t, "in memory ansible inventory for 'in memory ansible cli output'", inventory.Name())
}

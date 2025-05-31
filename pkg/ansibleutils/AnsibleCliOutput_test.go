package ansibleutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func Test_Constructor(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		output := ansibleutils.NewAnsibleCliOutput()
		require.EqualValues(t, "in memory ansible cli output", output.Name())
	})

	t.Run("number of hosts", func(t *testing.T) {
		output := ansibleutils.NewAnsibleCliOutput()
		require.EqualValues(t, 0, mustutils.Must(output.GetNumberOfHosts(ctx())))
	})
}

func Test_CreateInventory(t *testing.T) {
	output := ansibleutils.NewAnsibleCliOutput()
	inventory := output.CreateInventory()
	require.EqualValues(t, "in memory ansible inventory for 'in memory ansible cli output'", inventory.Name())
}

package ansibleutils_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/ansibleutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
)

func Test_AnsibleInventoryConstructor(t *testing.T) {
	require.EqualValues(
		t,
		"in memory ansible inventory",
		ansibleutils.NewAnsibleInventory().Name(),
	)
}

func Test_AddAndGetHostnames(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		inventory := ansibleutils.NewAnsibleInventory()

		hostnames, err := inventory.ListHostNames()
		require.NoError(t, err)
		require.EqualValues(t, []string{}, hostnames)
	})

	t.Run("One entry", func(t *testing.T) {
		inventory := ansibleutils.NewAnsibleInventory()

		_, err := inventory.CreateHostByName(ctx(), "test.example.net")
		require.NoError(t, err)

		hostnames, err := inventory.ListHostNames()
		require.NoError(t, err)
		require.EqualValues(t, []string{"test.example.net"}, hostnames)
	})

	t.Run("Two entries already sorted", func(t *testing.T) {
		inventory := ansibleutils.NewAnsibleInventory()

		_, err := inventory.CreateHostByName(ctx(), "a.example.net")
		require.NoError(t, err)

		_, err = inventory.CreateHostByName(ctx(), "test.example.net")
		require.NoError(t, err)

		hostnames, err := inventory.ListHostNames()
		require.NoError(t, err)
		require.EqualValues(t, []string{"a.example.net", "test.example.net"}, hostnames)
	})

	t.Run("Two entries", func(t *testing.T) {
		inventory := ansibleutils.NewAnsibleInventory()

		_, err := inventory.CreateHostByName(ctx(), "test.example.net")
		require.NoError(t, err)

		_, err = inventory.CreateHostByName(ctx(), "a.example.net")
		require.NoError(t, err)

		hostnames, err := inventory.ListHostNames()
		require.NoError(t, err)
		require.EqualValues(t, []string{"a.example.net", "test.example.net"}, hostnames)
	})
}

func Test_ParseInventoryJson(t *testing.T) {
	input := `{` + "\n"
	input += `	"_meta": {` + "\n"
	input += `		"hostvars": {}` + "\n"
	input += `	},` + "\n"
	input += `	"all": {` + "\n"
	input += `		"children": [` + "\n"
	input += `			"ungrouped"` + "\n"
	input += `		]` + "\n"
	input += `	},` + "\n"
	input += `	"ungrouped": {` + "\n"
	input += `		"hosts": [` + "\n"
	input += `			"one.example.net"` + "\n"
	input += `		]` + "\n"
	input += `	}` + "\n"
	input += `}` + "\n"

	inventory, err := ansibleutils.ParseInventoryJson(ctx(), input)
	require.NoError(t, err)

	require.EqualValues(t, 1, mustutils.Must(inventory.GetNumberOfHosts(ctx())))
	require.EqualValues(t, []string{"one.example.net"}, mustutils.Must(inventory.ListHostNames()))
	require.EqualValues(t, []string{"all", "ungrouped"}, mustutils.Must(inventory.ListGroupNames()))
}

func Test_CreateGroupByName(t *testing.T) {
	inventory := ansibleutils.NewAnsibleInventory()
	ctx := ctx()

	nGroups := 10

	for i := 0; i < nGroups; i++ {
		groupName := "group-" + strconv.Itoa(i)

		require.False(t, mustutils.Must(inventory.GroupByNameExists(ctx, groupName)))

		for k := 0; k < 2; k++ {
			_, err := inventory.CreateGroupByName(ctx, groupName)
			require.NoError(t, err)

			require.True(t, mustutils.Must(inventory.GroupByNameExists(ctx, groupName)))
		}
	}

	require.Len(t, mustutils.Must(inventory.ListGroupNames()), nGroups+1) // created groups + "all"

}

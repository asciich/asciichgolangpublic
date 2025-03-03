package ansibleutils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AnsibleInventoryConstructor(t *testing.T) {
	require.EqualValues(
		t,
		"in memory ansible inventory",
		NewAnsibleInventory().Name(),
	)
}

func Test_AddAndGetHostnames(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		inventory := NewAnsibleInventory()
		require.EqualValues(t, []string{}, inventory.MustListHostNames())
	})

	t.Run("One entry", func(t *testing.T) {
		inventory := NewAnsibleInventory()
		inventory.MustCreateHostByName(ctx(), "test.example.net")
		require.EqualValues(t, []string{"test.example.net"}, inventory.MustListHostNames())
	})

	t.Run("Two entries already sorted", func(t *testing.T) {
		inventory := NewAnsibleInventory()
		inventory.MustCreateHostByName(ctx(), "a.example.net")
		inventory.MustCreateHostByName(ctx(), "test.example.net")
		require.EqualValues(t, []string{"a.example.net", "test.example.net"}, inventory.MustListHostNames())
	})

	t.Run("Two entries", func(t *testing.T) {
		inventory := NewAnsibleInventory()
		inventory.MustCreateHostByName(ctx(), "test.example.net")
		inventory.MustCreateHostByName(ctx(), "a.example.net")
		require.EqualValues(t, []string{"a.example.net", "test.example.net"}, inventory.MustListHostNames())
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

	inventory := MustParseInventoryJson(ctx(), input)

	require.EqualValues(t, 1, inventory.MustGetNumberOfHosts(ctx()))
	require.EqualValues(t, []string{"one.example.net"}, inventory.MustListHostNames())
	require.EqualValues(t, []string{"all", "ungrouped"}, inventory.MustListGroupNames())
}

func Test_CreateGroupByName(t *testing.T) {
	inventory := NewAnsibleInventory()
	ctx := ctx()

	nGroups := 10

	for i := 0; i < nGroups; i++ {
		groupName := "group-" + strconv.Itoa(i)

		require.False(t, inventory.MustGroupByNameExists(ctx, groupName))

		for k := 0; k < 2; k++ {
			inventory.MustCreateGroupByName(ctx, groupName)
			require.True(t, inventory.MustGroupByNameExists(ctx, groupName))
		}
	}

	require.Len(
		t,
		inventory.MustListGroupNames(),
		nGroups+1, // created groups + "all"
	)
}

package ansibleutils

import (
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
		require.EqualValues(t, []string{}, inventory.GetHostNames())
	})

	t.Run("One entry", func(t *testing.T) {
		inventory := NewAnsibleInventory()
		inventory.MustAddHostByName(ctx(), "test.example.net")
		require.EqualValues(t, []string{"test.example.net"}, inventory.GetHostNames())
	})

	t.Run("Two entries already sorted", func(t *testing.T) {
		inventory := NewAnsibleInventory()
		inventory.MustAddHostByName(ctx(), "a.example.net")
		inventory.MustAddHostByName(ctx(), "test.example.net")
		require.EqualValues(t, []string{"a.example.net", "test.example.net"}, inventory.GetHostNames())
	})

	t.Run("Two entries", func(t *testing.T) {
		inventory := NewAnsibleInventory()
		inventory.MustAddHostByName(ctx(), "test.example.net")
		inventory.MustAddHostByName(ctx(), "a.example.net")
		require.EqualValues(t, []string{"a.example.net", "test.example.net"}, inventory.GetHostNames())
	})
}

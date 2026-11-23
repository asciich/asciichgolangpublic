package headscalegeneric_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/iputils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalegeneric"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
)

func getFilledNode(t *testing.T) *headscalegeneric.Node {
	ctx := context.Background()

	lastSeen := time.Date(2026, 11, 21, 12, 48, 50, 0, time.UTC)

	node := headscalegeneric.NewNode()
	require.NoError(t, node.SetID(1))
	require.NoError(t, node.SetHostname("headscale-ssh-admin"))
	require.NoError(t, node.SetName("headscale-ssh-admin"))
	require.NoError(t, node.SetMachineKey("cMIcn"))
	require.NoError(t, node.SetNodeKey("6MOQG"))
	require.NoError(t, node.SetUser("reto"))
	require.NoError(t, node.SetTags([]string{"tag:backup"}))
	require.NoError(t, node.SetIPAddresses(ctx, []string{"100.64.0.1", "fd7a:115c:a1e0::1"}))
	require.NoError(t, node.SetEphemeral(false))
	require.NoError(t, node.SetLastSeen(&lastSeen))
	require.NoError(t, node.SetExpiry(nil))
	require.NoError(t, node.SetOnline(true))

	return node
}

func TestNode_ImplementsInterface(t *testing.T) {
	t.Run("node implements the Node interface", func(t *testing.T) {
		var node headscaleinterfaces.Node = headscalegeneric.NewNode()
		require.NotNil(t, node)
	})
}

func TestNode_GettersOfEmptyNode(t *testing.T) {
	node := headscalegeneric.NewNode()

	t.Run("id", func(t *testing.T) {
		_, err := node.GetID()
		require.Error(t, err)
	})

	t.Run("hostname", func(t *testing.T) {
		_, err := node.GetHostname()
		require.Error(t, err)
	})

	t.Run("name", func(t *testing.T) {
		_, err := node.GetName()
		require.Error(t, err)
	})

	t.Run("machine key", func(t *testing.T) {
		_, err := node.GetMachineKey()
		require.Error(t, err)
	})

	t.Run("node key", func(t *testing.T) {
		_, err := node.GetNodeKey()
		require.Error(t, err)
	})

	t.Run("user", func(t *testing.T) {
		_, err := node.GetUser()
		require.Error(t, err)
	})

	t.Run("tags are empty and not nil", func(t *testing.T) {
		tags, err := node.GetTags()
		require.NoError(t, err)
		require.NotNil(t, tags)
		require.Empty(t, tags)
	})

	t.Run("ip addresses are empty and not nil", func(t *testing.T) {
		ips, err := node.GetIPAddresses()
		require.NoError(t, err)
		require.NotNil(t, ips)
		require.Empty(t, ips)
	})

	t.Run("not expired without expiry", func(t *testing.T) {
		expired, err := node.IsExpired()
		require.NoError(t, err)
		require.False(t, expired)
	})
}

func TestNode_SetAndGet(t *testing.T) {
	node := getFilledNode(t)

	t.Run("id", func(t *testing.T) {
		id, err := node.GetID()
		require.NoError(t, err)
		require.EqualValues(t, 1, id)
	})

	t.Run("hostname", func(t *testing.T) {
		hostname, err := node.GetHostname()
		require.NoError(t, err)
		require.Equal(t, "headscale-ssh-admin", hostname)
	})

	t.Run("user", func(t *testing.T) {
		user, err := node.GetUser()
		require.NoError(t, err)
		require.Equal(t, "reto", user)
	})

	t.Run("tags", func(t *testing.T) {
		tags, err := node.GetTags()
		require.NoError(t, err)
		require.Equal(t, []string{"tag:backup"}, tags)
	})

	t.Run("ip addresses", func(t *testing.T) {
		ips, err := node.GetIPAddresses()
		require.NoError(t, err)
		require.Equal(t, []string{"100.64.0.1", "fd7a:115c:a1e0::1"}, ips)
	})

	t.Run("online", func(t *testing.T) {
		online, err := node.IsOnline()
		require.NoError(t, err)
		require.True(t, online)
	})

	t.Run("ephemeral", func(t *testing.T) {
		ephemeral, err := node.IsEphemeral()
		require.NoError(t, err)
		require.False(t, ephemeral)
	})
}

func TestNode_SetInvalidValues(t *testing.T) {
	ctx := context.Background()
	node := headscalegeneric.NewNode()

	t.Run("id 0", func(t *testing.T) {
		require.Error(t, node.SetID(0))
	})

	t.Run("empty hostname", func(t *testing.T) {
		require.Error(t, node.SetHostname(""))
	})

	t.Run("empty name", func(t *testing.T) {
		require.Error(t, node.SetName(""))
	})

	t.Run("empty machine key", func(t *testing.T) {
		require.Error(t, node.SetMachineKey(""))
	})

	t.Run("empty node key", func(t *testing.T) {
		require.Error(t, node.SetNodeKey(""))
	})

	t.Run("empty user", func(t *testing.T) {
		require.Error(t, node.SetUser(""))
	})

	t.Run("empty tag", func(t *testing.T) {
		require.Error(t, node.SetTags([]string{"tag:ok", ""}))
	})

	t.Run("invalid ip address", func(t *testing.T) {
		err := node.SetIPAddresses(ctx, []string{"100.64.0.1", "not-an-ip"})
		require.Error(t, err)
		require.True(t, iputils.IsInvalidIPError(err))
	})

	t.Run("empty ip address", func(t *testing.T) {
		require.Error(t, node.SetIPAddresses(ctx, []string{""}))
	})

	t.Run("invalid ip addresses do not modify the node", func(t *testing.T) {
		ips, err := node.GetIPAddresses()
		require.NoError(t, err)
		require.Empty(t, ips)
	})
}

func TestNode_NilNode(t *testing.T) {
	var node *headscalegeneric.Node

	t.Run("getter", func(t *testing.T) {
		_, err := node.GetID()
		require.Error(t, err)
	})

	t.Run("setter", func(t *testing.T) {
		require.Error(t, node.SetOnline(true))
	})

	t.Run("deep copy", func(t *testing.T) {
		_, err := node.GetDeepCopy()
		require.Error(t, err)
	})
}

func TestNode_TimeHandling(t *testing.T) {
	t.Run("zero time expiry maps to nil", func(t *testing.T) {
		node := headscalegeneric.NewNode()

		zero := time.Time{} // 0001-01-01 00:00:00
		require.NoError(t, node.SetExpiry(&zero))

		expiry, err := node.GetExpiry()
		require.NoError(t, err)
		require.Nil(t, expiry)
	})

	t.Run("zero time last seen maps to nil", func(t *testing.T) {
		node := headscalegeneric.NewNode()

		zero := time.Time{}
		require.NoError(t, node.SetLastSeen(&zero))

		lastSeen, err := node.GetLastSeen()
		require.NoError(t, err)
		require.Nil(t, lastSeen)
	})

	t.Run("last seen is returned", func(t *testing.T) {
		node := getFilledNode(t)

		lastSeen, err := node.GetLastSeen()
		require.NoError(t, err)
		require.NotNil(t, lastSeen)
		require.True(t, lastSeen.Equal(time.Date(2026, 11, 21, 12, 48, 50, 0, time.UTC)))
	})

	t.Run("expiry in the past is expired", func(t *testing.T) {
		node := headscalegeneric.NewNode()

		past := time.Now().Add(-time.Hour)
		require.NoError(t, node.SetExpiry(&past))

		expired, err := node.IsExpired()
		require.NoError(t, err)
		require.True(t, expired)
	})

	t.Run("expiry in the future is not expired", func(t *testing.T) {
		node := headscalegeneric.NewNode()

		future := time.Now().Add(time.Hour)
		require.NoError(t, node.SetExpiry(&future))

		expired, err := node.IsExpired()
		require.NoError(t, err)
		require.False(t, expired)
	})
}

func TestNode_ReturnedSlicesAreCopies(t *testing.T) {
	t.Run("modifying returned tags does not change the node", func(t *testing.T) {
		node := getFilledNode(t)

		tags, err := node.GetTags()
		require.NoError(t, err)
		tags[0] = "modified"

		tagsAgain, err := node.GetTags()
		require.NoError(t, err)
		require.Equal(t, []string{"tag:backup"}, tagsAgain)
	})

	t.Run("modifying returned ip addresses does not change the node", func(t *testing.T) {
		node := getFilledNode(t)

		ips, err := node.GetIPAddresses()
		require.NoError(t, err)
		ips[0] = "modified"

		ipsAgain, err := node.GetIPAddresses()
		require.NoError(t, err)
		require.Equal(t, "100.64.0.1", ipsAgain[0])
	})
}

func TestNode_GetDeepCopy(t *testing.T) {
	ctx := context.Background()

	t.Run("copy has the same values", func(t *testing.T) {
		node := getFilledNode(t)

		cp, err := node.GetDeepCopy()
		require.NoError(t, err)
		require.NotSame(t, node, cp)
		require.Equal(t, node, cp)
	})

	t.Run("tags are deep copied", func(t *testing.T) {
		node := getFilledNode(t)

		cp, err := node.GetDeepCopy()
		require.NoError(t, err)

		require.NoError(t, cp.SetTags([]string{"tag:other"}))

		tags, err := node.GetTags()
		require.NoError(t, err)
		require.Equal(t, []string{"tag:backup"}, tags)
	})

	t.Run("ip addresses are deep copied", func(t *testing.T) {
		node := getFilledNode(t)

		cp, err := node.GetDeepCopy()
		require.NoError(t, err)

		require.NoError(t, cp.SetIPAddresses(ctx, []string{"100.64.0.99"}))

		ips, err := node.GetIPAddresses()
		require.NoError(t, err)
		require.Equal(t, []string{"100.64.0.1", "fd7a:115c:a1e0::1"}, ips)
	})

	t.Run("times are deep copied", func(t *testing.T) {
		node := getFilledNode(t)

		cp, err := node.GetDeepCopy()
		require.NoError(t, err)

		newLastSeen := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
		require.NoError(t, cp.SetLastSeen(&newLastSeen))

		lastSeen, err := node.GetLastSeen()
		require.NoError(t, err)
		require.True(t, lastSeen.Equal(time.Date(2026, 11, 21, 12, 48, 50, 0, time.UTC)))
	})

	t.Run("copy of an empty node keeps empty slices non nil", func(t *testing.T) {
		cp, err := headscalegeneric.NewNode().GetDeepCopy()
		require.NoError(t, err)

		tags, err := cp.GetTags()
		require.NoError(t, err)
		require.NotNil(t, tags)
	})
}

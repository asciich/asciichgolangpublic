package ansiblegalaxyutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansiblegalaxyutils"
)

func Test_Options_GetAnsiblePath(t *testing.T) {
	t.Run("nil fails", func(t *testing.T) {
		path, err := ansiblegalaxyutils.GetAnsiblePath(nil)
		require.Error(t, err)
		require.Empty(t, path)
	})

	t.Run("Empty struct", func(t *testing.T) {
		type Options struct{}

		path, err := ansiblegalaxyutils.GetAnsiblePath(&Options{})
		require.NoError(t, err)
		require.EqualValues(t, path, "ansible")
	})

	t.Run("AnsibleVirtualenvPath in struct but not set", func(t *testing.T) {
		type Options struct {
			AnsibleVirtualenvPath string
		}

		path, err := ansiblegalaxyutils.GetAnsiblePath(&Options{})
		require.NoError(t, err)
		require.EqualValues(t, path, "ansible")
	})

	t.Run("AnsibleVirtualenvPath in struct", func(t *testing.T) {
		type Options struct {
			AnsibleVirtualenvPath string
		}

		path, err := ansiblegalaxyutils.GetAnsiblePath(&Options{
			AnsibleVirtualenvPath: "/opt/test",
		})
		require.NoError(t, err)
		require.EqualValues(t, path, "/opt/test/bin/ansible")
	})
}

func Test_Options_GetAnsibleGalaxyPath(t *testing.T) {
	t.Run("nil fails", func(t *testing.T) {
		path, err := ansiblegalaxyutils.GetAnsibleGalaxyPath(nil)
		require.Error(t, err)
		require.Empty(t, path)
	})

	t.Run("Empty struct", func(t *testing.T) {
		type Options struct{}

		path, err := ansiblegalaxyutils.GetAnsibleGalaxyPath(&Options{})
		require.NoError(t, err)
		require.EqualValues(t, path, "ansible-galaxy")
	})

	t.Run("AnsibleVirtualenvPath in struct but not set", func(t *testing.T) {
		type Options struct {
			AnsibleVirtualenvPath string
		}

		path, err := ansiblegalaxyutils.GetAnsibleGalaxyPath(&Options{})
		require.NoError(t, err)
		require.EqualValues(t, path, "ansible-galaxy")
	})

	t.Run("AnsibleVirtualenvPath in struct", func(t *testing.T) {
		type Options struct {
			AnsibleVirtualenvPath string
		}

		path, err := ansiblegalaxyutils.GetAnsibleGalaxyPath(&Options{
			AnsibleVirtualenvPath: "/opt/test",
		})
		require.NoError(t, err)
		require.EqualValues(t, path, "/opt/test/bin/ansible-galaxy")
	})
}

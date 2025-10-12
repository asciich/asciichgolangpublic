package ansibleparemeteroptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansibleparemeteroptions"
)

func Test_Options_GetAnsiblePath(t *testing.T) {
	t.Run("nil fails", func(t *testing.T) {
		path, err := ansibleparemeteroptions.GetAnsiblePath(nil)
		require.Error(t, err)
		require.Empty(t, path)
	})

	t.Run("Empty struct", func(t *testing.T) {
		type Options struct{}

		path, err := ansibleparemeteroptions.GetAnsiblePath(&Options{})
		require.NoError(t, err)
		require.EqualValues(t, path, "ansible")
	})

	t.Run("AnsibleVirtualenvPath in struct but not set", func(t *testing.T) {
		type Options struct {
			AnsibleVirtualenvPath string
		}

		path, err := ansibleparemeteroptions.GetAnsiblePath(&Options{})
		require.NoError(t, err)
		require.EqualValues(t, path, "ansible")
	})

	t.Run("AnsibleVirtualenvPath in struct", func(t *testing.T) {
		type Options struct {
			AnsibleVirtualenvPath string
		}

		path, err := ansibleparemeteroptions.GetAnsiblePath(&Options{
			AnsibleVirtualenvPath: "/opt/test",
		})
		require.NoError(t, err)
		require.EqualValues(t, path, "/opt/test/bin/ansible")
	})
}

func Test_Options_GetAnsibleGalaxyPath(t *testing.T) {
	t.Run("nil fails", func(t *testing.T) {
		path, err := ansibleparemeteroptions.GetAnsibleGalaxyPath(nil)
		require.Error(t, err)
		require.Empty(t, path)
	})

	t.Run("Empty struct", func(t *testing.T) {
		type Options struct{}

		path, err := ansibleparemeteroptions.GetAnsibleGalaxyPath(&Options{})
		require.NoError(t, err)
		require.EqualValues(t, path, "ansible-galaxy")
	})

	t.Run("AnsibleVirtualenvPath in struct but not set", func(t *testing.T) {
		type Options struct {
			AnsibleVirtualenvPath string
		}

		path, err := ansibleparemeteroptions.GetAnsibleGalaxyPath(&Options{})
		require.NoError(t, err)
		require.EqualValues(t, path, "ansible-galaxy")
	})

	t.Run("AnsibleVirtualenvPath in struct", func(t *testing.T) {
		type Options struct {
			AnsibleVirtualenvPath string
		}

		path, err := ansibleparemeteroptions.GetAnsibleGalaxyPath(&Options{
			AnsibleVirtualenvPath: "/opt/test",
		})
		require.NoError(t, err)
		require.EqualValues(t, path, "/opt/test/bin/ansible-galaxy")
	})
}


func Test_Options_GetAnsiblePlaybookPath(t *testing.T) {
	t.Run("nil fails", func(t *testing.T) {
		path, err := ansibleparemeteroptions.GetAnsiblePlaybookPath(nil)
		require.Error(t, err)
		require.Empty(t, path)
	})

	t.Run("Empty struct", func(t *testing.T) {
		type Options struct{}

		path, err := ansibleparemeteroptions.GetAnsiblePlaybookPath(&Options{})
		require.NoError(t, err)
		require.EqualValues(t, path, "ansible-playbook")
	})

	t.Run("AnsibleVirtualenvPath in struct but not set", func(t *testing.T) {
		type Options struct {
			AnsibleVirtualenvPath string
		}

		path, err := ansibleparemeteroptions.GetAnsiblePlaybookPath(&Options{})
		require.NoError(t, err)
		require.EqualValues(t, path, "ansible-playbook")
	})

	t.Run("AnsibleVirtualenvPath in struct", func(t *testing.T) {
		type Options struct {
			AnsibleVirtualenvPath string
		}

		path, err := ansibleparemeteroptions.GetAnsiblePlaybookPath(&Options{
			AnsibleVirtualenvPath: "/opt/test",
		})
		require.NoError(t, err)
		require.EqualValues(t, path, "/opt/test/bin/ansible-playbook")
	})
}

package ansiblegalaxyutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansiblegalaxyutils"
)

func Test_ListInstalledCollectionsOptions_GetAnsiblePath(t *testing.T) {
	t.Run("Not set results in 'ansible'", func(t *testing.T) {
		options := &ansiblegalaxyutils.InstallCollectionOptions{}
		ansiblePath, err := options.GetAnsiblePath()
		require.NoError(t, err)
		require.EqualValues(t, "ansible", ansiblePath)
	})

	t.Run("Set virtualenv", func(t *testing.T) {
		options := &ansiblegalaxyutils.InstallCollectionOptions{
			AnsibleVirtualenvPath: "/opt/ve_test/",
		}
		ansiblePath, err := options.GetAnsiblePath()
		require.NoError(t, err)
		require.EqualValues(t, "/opt/ve_test/bin/ansible", ansiblePath)
	})
}

func Test_ListInstalledCollectionsOptions_GetAnsibleGalaxyPath(t *testing.T) {
	t.Run("Not set results in 'ansible-galaxy'", func(t *testing.T) {
		options := &ansiblegalaxyutils.InstallCollectionOptions{}
		ansiblePath, err := options.GetAnsibleGalaxyPath()
		require.NoError(t, err)
		require.EqualValues(t, "ansible-galaxy", ansiblePath)
	})

	t.Run("Set virtualenv", func(t *testing.T) {
		options := &ansiblegalaxyutils.InstallCollectionOptions{
			AnsibleVirtualenvPath: "/opt/ve_test/",
		}
		ansiblePath, err := options.GetAnsibleGalaxyPath()
		require.NoError(t, err)
		require.EqualValues(t, "/opt/ve_test/bin/ansible-galaxy", ansiblePath)
	})
}

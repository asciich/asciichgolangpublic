package filesoptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
)

func Test_GetPermissionsString(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		options := filesoptions.ChmodOptions{}
		
		permissionsString, err := options.GetPermissionsString()
		require.Error(t, err)
		require.Empty(t, permissionsString)
	})

	t.Run("u=r", func(t *testing.T) {
		options := filesoptions.ChmodOptions{
			PermissionsString: "u=r",
		}

		permissionsString, err := options.GetPermissionsString()
		require.NoError(t, err)
		require.EqualValues(t, "u=r", permissionsString)
	})
} 

func Test_GetPermissions(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		options := filesoptions.ChmodOptions{}
		
		permissions, err := options.GetPermissions()
		require.Error(t, err)
		require.Zero(t, permissions)
	})

	t.Run("u=r", func(t *testing.T) {
		options := filesoptions.ChmodOptions{
			PermissionsString: "u=r",
		}

		permissions, err := options.GetPermissions()
		require.NoError(t, err)
		require.EqualValues(t, 0400, permissions)
	})
}

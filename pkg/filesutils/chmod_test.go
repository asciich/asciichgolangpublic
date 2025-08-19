package filesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/unixfilepermissionsutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestFile_Chmod(t *testing.T) {
	t.Run("commandexecutor", func(t *testing.T) {
		ctx := getCtx()
		pathToTest, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)

		executor := commandexecutorbashoo.Bash()
		permissions, err := commandexecutorfile.GetAccessPermissionsString(executor, pathToTest)
		require.NoError(t, err)
		require.NotEqualValues(t, "u=rw,g=r,o=r", permissions)

		for range 2 {
			err = commandexecutorfile.Chmod(ctx, executor, pathToTest, &filesoptions.ChmodOptions{
				PermissionsString: "u=rw,g=r,o=r",
			})
			require.NoError(t, err)

			permissions, err = commandexecutorfile.GetAccessPermissionsString(executor, pathToTest)
			require.NoError(t, err)
			require.EqualValues(t, "u=rw,g=r,o=r", permissions)
		}
	})

	t.Run("nativefiles", func(t *testing.T) {
		ctx := getCtx()
		pathToTest, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)

		permissions, err := nativefiles.GetAccessPermissionsString(pathToTest)
		require.NoError(t, err)
		require.NotEqualValues(t, "u=rw,g=r,o=r", permissions)

		for range 2 {
			err = nativefiles.Chmod(ctx, pathToTest, &filesoptions.ChmodOptions{
				PermissionsString: "u=rw,g=r,o=r",
			})
			require.NoError(t, err)

			permissions, err = nativefiles.GetAccessPermissionsString(pathToTest)
			require.NoError(t, err)
			require.EqualValues(t, "u=rw,g=r,o=r", permissions)
		}
	})
}

func TestFile_ChmodObjectOriented(t *testing.T) {
	tests := []struct {
		implementationName       string
		permissionsString        string
		expectedPermissionString string
	}{
		{"localFile", "u=rw,g=,o=", "u=rw,g=,o="},
		{"localCommandExecutorFile", "u=rw,g=,o=", "u=rw,g=,o="},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				toTest := getTemporaryFileToTest(tt.implementationName)
				defer toTest.Delete(ctx, &filesoptions.DeleteOptions{})

				err := toTest.Chmod(
					ctx,
					&filesoptions.ChmodOptions{
						PermissionsString: tt.permissionsString,
					},
				)
				require.NoError(t, err)

				accessPermissionsString, err := toTest.GetAccessPermissionsString()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedPermissionString, accessPermissionsString)

				accessPermissions, err := toTest.GetAccessPermissions()
				require.NoError(t, err)
				require.EqualValues(t, unixfilepermissionsutils.MustGetPermissionsValue(tt.expectedPermissionString), accessPermissions)
			},
		)
	}
}

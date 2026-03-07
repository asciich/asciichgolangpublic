package filesutils_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
)

// To run this test use:
//
//	bash -c "RUN_SUDO_TEST=1 go test -v $(git rev-parse --show-toplevel)/pkg/filesutils -run Test_CreateFileUsingSudo"
func Test_CreateFileUsingSudo(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localFile"},
		{"localCommandExecutorFile"},
		{"commandExecutorFileExec"},
		{"commandExecutorFileBash"},
	}

	for _, tt := range tests {
		t.Run("no root permission denied", func(t *testing.T) {
			ctx := getCtx()

			const testPath = "/testfile"

			// Creating the test file in the root directory without sudo failed:
			sourceFile := getFileToTest(tt.implementationName, testPath)

			// Hint: Ensure the /testfile is absent, otherwise this test failes.
			// The idempotent written Create function will skip the attempt to create the file if it already exists.
			err := sourceFile.Create(ctx, &filesoptions.CreateOptions{})
			require.Error(t, err)

			require.Contains(t, strings.ToLower(err.Error()), "permission denied")
		})
	}

	for _, tt := range tests {
		t.Run("with root permission granted", func(t *testing.T) {
			const envName = "RUN_SUDO_TEST"
			if os.Getenv(envName) != "1" {
				t.Skipf("Sudo tests skipped since '%s' not set.'", envName)
			}

			ctx := getCtx()

			sourceFile := getFileToTest(tt.implementationName, "/testfile")
			defer func() {
				err := sourceFile.Delete(ctx, &filesoptions.DeleteOptions{UseSudo: true})
				require.NoError(t, err)
			}()
			err := sourceFile.Create(ctx, &filesoptions.CreateOptions{UseSudo: true})
			require.NoError(t, err)
		})
	}

	for _, tt := range tests {
		for _, permissionString := range []string{"u=rwx,g=r,o="} {
			t.Run("chmod "+tt.implementationName+" "+permissionString, func(t *testing.T) {
				ctx := getCtx()

				testfile := getTemporaryFileToTest(tt.implementationName)
				defer func() {
					err := testfile.Delete(ctx, &filesoptions.DeleteOptions{})
					require.NoError(t, err)
				}()

				exists, err := testfile.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				err = testfile.Chmod(ctx, &filesoptions.ChmodOptions{
					PermissionsString: permissionString,
				})
				require.NoError(t, err)

				got, err := testfile.GetAccessPermissionsString()
				require.NoError(t, err)

				require.EqualValues(t, permissionString, got)
			})
		}
	}
}

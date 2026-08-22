package filesutils_test

import (
	"os/user"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// TestFile_Chown tests Chown method
func TestFile_Chown(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativefilesoo"},
		{"commandExecutorFileExec"},
		{"commandExecutorFileBash"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				fileToTest := getTemporaryFileToTest(tt.implementationName)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				// Get current user
				currentUser, err := user.Current()
				require.NoError(t, err)

				// Chown to current user (should succeed)
				err = fileToTest.Chown(ctx, &parameteroptions.ChownOptions{
					UserName: currentUser.Username,
				})
				require.NoError(t, err)
			},
		)
	}
}

package filesutils_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
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
			defer mustutils.Must0(sourceFile.Delete(ctx, &filesoptions.DeleteOptions{UseSudo: true}))
			err := sourceFile.Create(ctx, &filesoptions.CreateOptions{UseSudo: true})
			require.NoError(t, err)
		})
	}
}

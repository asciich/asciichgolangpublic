package filesutils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func getDirectoryToTest(implementationName string, testPath string) (directory filesinterfaces.Directory) {
	if implementationName == "localDirectory" {
		ctx := contextutils.ContextVerbose()
		dir, err := files.GetLocalDirectoryByPath(ctx, testPath)
		if err != nil {
			logging.LogGoErrorFatal(err)
		}

		return dir
	}

	if implementationName == "localCommandExecutorDirectory" {
		return mustutils.Must(files.GetLocalCommandExecutorDirectoryByPath(testPath))
	}

	if implementationName == "nativedirectoryoo" {
		return mustutils.Must(nativefilesoo.NewDirectoryByPath(testPath))
	}

	panic(fmt.Sprintf("unknown implementationName='%s'", implementationName))
}

func TestGetLocalPath(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localDirectory"},
		{"localCommandExecutorDirectory"},
		{"nativedirectoryoo"},
	}

	for _, tt := range tests {
		t.Run("getlocalpath_"+tt.implementationName, func(t *testing.T) {
			const testPath = "/testfile"

			sourceFile := getDirectoryToTest(tt.implementationName, testPath)

			localPath, err := sourceFile.GetLocalPath()
			require.NoError(t, err)

			require.EqualValues(t, "/testfile", localPath)
		})
	}
}

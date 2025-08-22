package datatypes_test

// Tests in this file are used to validate design decisions.

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func getDatatypesPath(t *testing.T) string {
	_, thisFilePath, _, ok := runtime.Caller(0)
	require.True(t, ok)

	return filepath.Dir(thisFilePath)

}

// This is kind of a self test to validate getting this package path correctly:
func Test_getDatatypesPath(t *testing.T) {
	t.Run("is absolute path", func(t *testing.T) {
		path := getDatatypesPath(t)
		require.True(t, pathsutils.IsAbsolutePath(path))
	})

	t.Run("is directory", func(t *testing.T) {
		path := getDatatypesPath(t)
		require.True(t, nativefiles.IsDir(getCtx(), path))
	})

	t.Run("Ends with datatype", func(t *testing.T) {
		path := getDatatypesPath(t)
		require.True(t, strings.HasSuffix(path, "/datatypes"))
	})
}

// Validates no fmt.Errorf are used/ returned in this package and it's sub packages.
// TracedErrors should be used instead.
func Test_DatatypesReturnsNoFmtErrorf(t *testing.T) {
	ctx := getCtx()
	goFiles, err := nativefiles.ListFiles(ctx, getDatatypesPath(t), &parameteroptions.ListFileOptions{})
	require.NoError(t, err)

	const searchString = "fmt." + "Errorf"

	regexOkString := regexp.MustCompile(`"` + searchString + `"`)

	for _, sourcePath := range goFiles {
		if filepath.Base(sourcePath) == "gettypename.go" {
			// gettypename.go is the only exception to avoid cyclic import.
			continue
		}

		t.Run(sourcePath, func(t *testing.T) {
			content, err := os.ReadFile(sourcePath)
			require.NoError(t, err)

			// writing fmt.Errorf in a comment is ok, so comments are removed.
			content = []byte(stringsutils.RemoveComments(string(content)))

			content = regexOkString.ReplaceAll(content, []byte("***"))

			require.Falsef(t, strings.Contains(string(content), searchString), "Found '%s' in '%s'. Should use tracederrors.TracedErrorf instead.", searchString, sourcePath)
		})
	}
}

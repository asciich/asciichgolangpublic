package files_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// Use t.NoError(t, err) to check for errors. Do not use less specific t.Nil(t, err)
func Test_UseTNoErrorToCheckForNoError(t *testing.T) {
	ctx := getCtx()
	sourceFiles, err := nativefiles.ListFiles(ctx, ".", &parameteroptions.ListFileOptions{})
	require.NoError(t, err)

	const absentPattern = "require.Nil(t, " + " " + "err)"

	for _, f := range sourceFiles {
		contains, err := nativefiles.Contains(ctx, f, absentPattern)
		require.NoError(t, err)

		require.Falsef(t, contains, "Found obsolete '%s' in '%s'.", absentPattern, f)
	}
}

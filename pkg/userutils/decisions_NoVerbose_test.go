package userutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// the legacy `verbose` boolean has to be replaced by ctx.
func Test_VerboseAbsent(t *testing.T) {
	ctx := getCtx()
	sourceFiles, err := nativefiles.ListFiles(ctx, ".", &parameteroptions.ListFileOptions{})
	require.NoError(t, err)

	const absentPattern = "verbose" + " " + "bool"

	for _, f := range sourceFiles {
		contains, err := nativefiles.Contains(ctx, f, absentPattern)
		require.NoError(t, err)

		require.Falsef(t, contains, "Found obsolete '%s' in '%s'.", absentPattern, f)
	}
}
